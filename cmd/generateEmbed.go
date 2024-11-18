package cmd

func generateEmbedGo() string {
	return `package main

import "embed"

//go:embed assets/*
var embeddedFiles embed.FS
`
}
