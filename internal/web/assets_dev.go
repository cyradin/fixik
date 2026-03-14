//go:build !prod

package web

import (
	"io/fs"
	"os"
)

func getStaticFS() fs.FS {
	fsys := os.DirFS("internal/web")
	sub, err := fs.Sub(fsys, "static")
	if err != nil {
		panic(err)
	}

	return sub
}
