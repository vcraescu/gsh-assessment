package web

import (
	"embed"
	_ "embed"
	"io/fs"
)

var (
	//go:embed static/*
	staticFS embed.FS
	FS, _    = fs.Sub(staticFS, "static")
)
