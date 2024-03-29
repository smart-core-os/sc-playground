package ui

// Package embed exposes our uis so they can be embedded inside go applications.
// This is not public API.
import (
	"embed"
	"io/fs"
)

//go:embed playground/dist
var playground embed.FS
var Playground, _ = fs.Sub(playground, "playground/dist/embed")
var NotFound, _ = fs.ReadFile(playground, "playground/dist/404.html")
