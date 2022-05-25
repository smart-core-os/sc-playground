package ui

// Package embed exposes our uis so they can be embedded inside go applications.
// This is not public API.
import (
	"embed"
	"io/fs"
)

//go:embed all:playground/dist
var playground embed.FS
var Playground, _ = fs.Sub(playground, "playground/dist")
