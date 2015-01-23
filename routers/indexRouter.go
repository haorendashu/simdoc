package routers

import (
	"github.com/flosch/pongo2"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, m pongo2.Context) {
	indexTpl := pongo2.Must(pongo2.FromFile("views/index.html"))
	indexTpl.ExecuteWriter(m, w)
}
