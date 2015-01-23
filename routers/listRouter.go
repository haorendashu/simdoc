package routers

import (
	"github.com/flosch/pongo2"
	"net/http"
	"github.com/haorendashu/simdoc/sdcomm"
	"path/filepath"
	"os"
)

func List(w http.ResponseWriter, r *http.Request, m pongo2.Context) {
	listPath(w, m, sdcomm.RootPath)
}

func listPath(w http.ResponseWriter, m pongo2.Context, rootPath string) {
	listTpl := pongo2.Must(pongo2.FromFile("views/list.html"))
	paths := []string{}
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) (error) {
			path, err = filepath.Rel(sdcomm.RootPath, path)
			if err != nil {
				sdcomm.Logger.Printf("walk error: %v.\n", err)
				return nil
			}
			if path == "." {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
	m["paths"] = paths

	listTpl.ExecuteWriter(m, w)
}
