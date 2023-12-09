package web

import (
	"embed"
	_ "embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed static/*
	staticFS     embed.FS
	contentFS, _ = fs.Sub(staticFS, "static")
)

func StaticHandler(w http.ResponseWriter, r *http.Request) error {
	http.FileServer(http.FS(contentFS)).ServeHTTP(w, r)

	return nil
}
