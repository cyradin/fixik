//go:build prod

package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var staticFS embed.FS

// getStaticFS returns the embedded static filesystem for production
func getStaticFS() fs.FS {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}

	return sub
}
