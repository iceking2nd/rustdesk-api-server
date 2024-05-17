package frontend

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static
var StaticFS embed.FS

//go:embed static/icons
var IconsFS embed.FS

//go:embed static/scripts
var ScriptsFS embed.FS

//go:embed static/styles
var StylesFS embed.FS

func GetFS(f embed.FS, prefix string) http.FileSystem {
	nf, err := fs.Sub(f, fmt.Sprintf("static/%s", prefix))
	if err != nil {
		panic(err.Error())
	}
	return http.FS(nf)
}
