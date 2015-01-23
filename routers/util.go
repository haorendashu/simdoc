package routers

import (
	"net/http"
	"strings"
	"path/filepath"
	"github.com/flosch/pongo2"
)

func IsMarkdownFile(name string) bool {
	name = strings.ToLower(name)
	switch filepath.Ext(name) {
	case ".md", ".markdown", ".mdown", ".mkd":
		return true
	}
	return false
}

func IsTextFile(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	if strings.Index(contentType, "text/") != -1 {
		return contentType, true
	}
	return contentType, false
}

func IsImageFile(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	if strings.Index(contentType, "image/") != -1 {
		return contentType, true
	}
	return contentType, false
}

func Return404(w http.ResponseWriter) {
	NotFoundTpl := pongo2.Must(pongo2.FromFile("views/404.html"))
	m := make(map[string]interface{})
	NotFoundTpl.ExecuteWriter(m, w)
}

func Return500(w http.ResponseWriter) {
	ErrorTpl    := pongo2.Must(pongo2.FromFile("views/500.html"))
	m := make(map[string]interface{})
	ErrorTpl.ExecuteWriter(m, w)
}
