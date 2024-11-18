package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generatePlayerUpdateGo(importPath, entityName string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)
	// Convert the entity name to title case for the struct name
	c := cases.Title(language.Und)
	titleEntityName := c.String(entityName)
	// Get the first letter of the entity name in lowercase for the receiver
	receiver := strings.ToLower(string(entityName[0]))

	return fmt.Sprintf(`package %s

import (
    "math"

    "%s/settings"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

// handleInput handles the %s's input.
func (%s *%s) handleInput() {
    if ebiten.IsKeyPressed(ebiten.KeyRight) {
        %s.direction = Right
        %s.SetState(Running) // hardcoded state, change it to your own state
    }

    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        %s.direction = Left
        %s.SetState(Running)  // hardcoded state, change it to your own state
    }

    if !ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft) {
        %s.SetState(Idle)  // hardcoded state, change it to your own state
    }
}

// Update updates the %s's state based on input.
func (%s *%s) Update() {
    %s.handleInput()

    friction := 0.5
    accel := 0.5 + friction
    maxSpeed := 4.0
    jumpSpd := 10.0
    gravity := 0.75

    %s.Speed.Y += gravity

    // This means, %s is sliding on the wall, and we want to limit his vertical speed to 1.
    // You need to set the state of %s into 'sliding'.
    if %s.SlidingOnWall != nil && %s.Speed.Y > 1 {
        %s.Speed.Y = 1
    }

    // Horizontal movement is only possible when not wallsliding.
    if %s.SlidingOnWall == nil {
        if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxisValue(0, 0) > 0.1 {
            %s.Speed.X += accel
            %s.direction = Right
            %s.FacingRight = true
        }

        if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxisValue(0, 0) < -0.1 {
            %s.Speed.X -= accel
            %s.direction = Left
            %s.FacingRight = false
        }
    }

    // Apply friction and horizontal speed limiting.
    if %s.Speed.X > friction {
        %s.Speed.X -= friction
    } else if %s.Speed.X < -friction {
        %s.Speed.X += friction
    } else {
        %s.Speed.X = 0
    }

    // Prevent %s from moving too fast.
    if %s.Speed.X > maxSpeed {
        %s.Speed.X = maxSpeed
    } else if %s.Speed.X < -maxSpeed {
        %s.Speed.X = -maxSpeed
    }

    // Check for jumping.
    if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsGamepadButtonJustPressed(0, 0) {
        // If the player is on a platform, and is pressing down, then they should drop down through the platform.
        if (ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.GamepadAxisValue(0, 1) > 0.1) && %s.OnGround != nil && %s.OnGround.HasTags("platform") {
            %s.IgnorePlatform = %s.OnGround // the platform is ignored, so the player can drop through it due to gravity.
        } else {
            // If the player is not pressing down or the ground object does not have the "platform" tag, then they should jump normally.

            // If the player is on the ground, then they can jump.
            if %s.OnGround != nil {
                %s.Speed.Y = -jumpSpd
            } else if %s.SlidingOnWall != nil {
                // If the player is not on the ground, but is sliding on a wall, then they can walljump.
                %s.Speed.Y = -jumpSpd

                // If the wall is to the right of %s, the jump to the left, and vice versa.
                if %s.SlidingOnWall.Position.X > %s.Object.Position.X {
                    %s.Speed.X = -4
                } else {
                    %s.Speed.X = 4
                }

                // %s is no longer wallsliding.
                %s.SlidingOnWall = nil
            }
        }
    }

    // We handle horizontal movement separately from vertical movement. This is, conceptually, decomposing movement into two phases / axes.
    // By decomposing movement in this manner, we can handle each case properly (i.e. stop movement horizontally separately from vertical movement, as
    // necesseary). More can be seen on this topic over on this blog post on higherorderfun.com:
    // http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

    // dx is the horizontal delta movement variable (which is the Player's horizontal speed). If we come into contact with something, then it will
    // be that movement instead.
    dx := %s.Speed.X

    // Moving horizontally is done fairly simply; we just check to see if something solid is in front of us. If so, we move into contact with it
    // and stop horizontal movement speed. If not, then we can just move forward.

    if check := %s.Object.Check(%s.Speed.X, 0, "solid"); check != nil {
        // If we come into contact with a solid object, we move as close as possible to contact with it, and stop horizontal movement.
        dx = check.ContactWithCell(check.Cells[0]).X
        %s.Speed.X = 0

        // If you're in the air, then colliding with a wall object makes you start wall sliding.
        if %s.OnGround == nil {
            %s.SlidingOnWall = check.Objects[0]
        }
    }

    // Then we just apply the horizontal movement to the Player's Object. Easy-peasy.
    %s.Object.Position.X += dx

    // Now for the vertical movement; it's the most complicated because we can land on different types of objects and need
    // to treat them all differently, but overall, it's not bad.

    // First, we set OnGround to be nil, in case we don't end up standing on anything.
    %s.OnGround = nil

    // dy is the delta movement downward, and is the vertical movement by default; similarly to dx, if we come into contact with
    // something, this will be changed to move to contact instead.

    dy := %s.Speed.Y

    // We want to be sure to lock vertical movement to a maximum of the size of the Cells within the Space
    // so we don't miss any collisions by tunneling through.
    dy = math.Max(math.Min(dy, settings.MaxCellSize), -settings.MaxCellSize)

    // We're going to check for collision using dy (which is vertical movement speed), but add one when moving downwards to look a bit deeper down
    // into the ground for solid objects to land on, specifically.
    checkDistance := dy
    if dy >= 0 {
        checkDistance++
    }

    // We check for any solid / stand-able objects. In actuality, there aren't any other Objects
    // with other tags in this Space, so we don't -have- to specify any tags, but it's good to be specific for clarity in this example.
    if check := %s.Object.Check(0, checkDistance, "solid", "platform", "ramp"); check != nil {

        // So! Firstly, we want to see if we jumped up into something that we can slide around horizontally to avoid bumping the Player's head.

        // Sliding around a misspaced jump is a small thing that makes jumping a bit more forgiving, and is something different polished platformers
        // (like the 2D Mario games) do to make it a smidge more comfortable to play. For a visual example of this, see this excellent devlog post
        // from the extremely impressive indie game, Leilani's Island: https://forums.tigsource.com/index.php?topic=46289.msg1387138#msg1387138

        // To accomplish this sliding, we simply call Collision.SlideAgainstCell() to see if we can slide.
        // We pass the first cell, and tags that we want to avoid when sliding (i.e. we don't want to slide into cells that contain other solid objects).

        slide, slideOK := check.SlideAgainstCell(check.Cells[0], "solid")

        // We further ensure that we only slide if:
        // 1) We're jumping up into something (dy < 0),
        // 2) If the cell we're bumping up against contains a solid object,
        // 3) If there was, indeed, a valid slide left or right, and
        // 4) If the proposed slide is less than 8 pixels in horizontal distance. (This is a relatively arbitrary number that just so happens to be half the
        // width of a cell. This is to ensure the player doesn't slide too far horizontally.)

        if dy < 0 && check.Cells[0].ContainsTags("solid") && slideOK && math.Abs(slide.X) <= settings.MaxCellSize/2 {

            // If we are able to slide here, we do so. No contact was made, and vertical speed (dy) is maintained upwards.
            %s.Object.Position.X += slide.X

        } else {

            // If sliding -fails-, that means the Player is jumping directly onto or into something, and we need to do more to see if we need to come into
            // contact with it. Let's press on!

            // First, we check for ramps. For ramps, we can't simply check for collision with Check(), as that's not precise enough. We need to get a bit
            // more information, and so will do so by checking its Shape (a triangular ConvexPolygon, as defined in WorldPlatformer.Init()) against the
            // Player's Shape (which is also a rectangular ConvexPolygon).

            // We get the ramp by simply filtering out Objects with the "ramp" tag out of the objects returned in our broad Check(), and grabbing the first one
            // if there's any at all.
            if ramps := check.ObjectsByTags("ramp"); len(ramps) > 0 {

                // For simplicity, this code assumes we can only stand on one ramp at a time as there is only one ramp in this example.
                // This is exemplified by the ramp := ramps[0] line.
                // In actuality, if there was a possibility to have a potential collision with multiple ramps (i.e. a ramp that sits on another ramp, and the player running down
                // one onto the other), the collision testing code should probably go with the ramp with the highest confirmed intersection point out of the two.

                ramp := ramps[0]

                // Next, we see if there's been an intersection between the two Shapes using Shape.Intersection. We pass the ramp's shape, and also the movement
                // we're trying to make horizontally, as this makes the Intersection function return the next y-position while moving, not the one directly
                // underneath the %s. This would keep the player from getting "stuck" when walking up a ramp into the top of a solid block, if there weren't
                // a landing at the top and bottom of the ramp.

                // We use 8 here for the Y-delta so that we can easily see if you're running down the ramp (in which case you're probably in the air as you
                // move faster than you can fall in this example). This way we can maintain contact so you can always jump while running down a ramp. We only
                // continue with coming into contact with the ramp as long as you're not moving upwards (i.e. jumping).

                if contactSet := %s.Object.Shape.Intersection(dx, settings.MaxCellSize/2, ramp.Shape); dy >= 0 && contactSet != nil {

                    // If Intersection() is successful, a ContactSet is returned. A ContactSet contains information regarding where
                    // two Shapes intersect, like the individual points of contact, the center of the contacts, and the MTV, or
                    // Minimum Translation Vector, to move out of contact.

                    // Here, we use ContactSet.TopmostPoint() to get the top-most contact point as an indicator of where
                    // we want the player's feet to be. Then we just set that position with a tiny bit of collision margin,
                    // and we're done.

                    dy = contactSet.TopmostPoint().Y - %s.Object.Bottom() + 0.1
                    %s.OnGround = ramp
                    %s.Speed.Y = 0

                }

            }

            // Platforms are next; here, we just see if the platform is not being ignored by attempting to drop down,
            // if the player is falling on the platform (as otherwise he would be jumping through platforms), and if the platform is low enough
            // to land on. If so, we stand on it.

            // Because there's a moving floating platform, we use Collision.ContactWithObject() to ensure the player comes into contact
            // with the top of the platform object. An alternative would be to use Collision.ContactWithCell(), but that would be only if the
            // platform didn't move and were aligned with the Space's grid.

            if platforms := check.ObjectsByTags("platform"); len(platforms) > 0 {

                platform := platforms[0]

                if platform != %s.IgnorePlatform && dy >= 0 && %s.Object.Bottom() < platform.Position.Y+4 {
                    dy = check.ContactWithObject(platform).Y
                    %s.OnGround = platform
                    %s.Speed.Y = 0
                }

            }

            // Finally, we check for simple solid ground. If we haven't had any success in landing previously, or the solid ground
            // is higher than the existing ground (like if the platform passes underneath the ground, or we're walking off of solid ground
            // onto a ramp), we stand on it instead. We don't check for solid collision first because we want any ramps to override solid
            // ground (so that you can walk onto the ramp, rather than sticking to solid ground).

            // We use ContactWithObject() here because otherwise, we might come into contact with the moving platform's cells (which, naturally,
            // would be selected by a Collision.ContactWithCell() call because the cell is closest to the Player).

            if solids := check.ObjectsByTags("solid"); len(solids) > 0 && (%s.OnGround == nil || %s.OnGround.Position.Y >= solids[0].Position.Y) {
                dy = check.ContactWithObject(solids[0]).Y
                %s.Speed.Y = 0

                // We're only on the ground if we land on it (if the object's Y is greater than the player's).
                if solids[0].Position.Y > %s.Object.Position.Y {
                    %s.OnGround = solids[0]
                }

            }

            if %s.OnGround != nil {
                %s.SlidingOnWall = nil  // Player's on the ground, so no wallsliding anymore.
                %s.IgnorePlatform = nil // Player's on the ground, so reset which platform is being ignored.
            }

        }

    }

    // Move the object on dy.
    %s.Object.Position.Y += dy

    wallNext := 1.0
    if !%s.FacingRight {
        wallNext = -1
    }

    // If the wall next to the Player runs out, stop wall sliding.
    if c := %s.Object.Check(wallNext, 0, "solid"); %s.SlidingOnWall != nil && c == nil {
        %s.SlidingOnWall = nil
    }

    %s.Object.Update() // Update the player's position in the space.

    // And that's it! Now move the sprite to the object's position.
    %s.x = %s.Object.Position.X
    %s.y = %s.Object.Position.Y
}

func (%s *%s) GetX() float64 {
    return %s.x
}

func (%s *%s) GetY() float64 {
    return %s.y
}
`, lowerEntityName, importPath, lowerEntityName, receiver, titleEntityName, receiver, receiver, receiver, receiver, receiver, lowerEntityName, receiver, titleEntityName, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, titleEntityName, receiver, receiver, titleEntityName, receiver)
}
