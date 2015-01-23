package routers

import (
	"github.com/flosch/pongo2"
	"net/http"
	"github.com/russross/blackfriday"
	"os"
	"path/filepath"
	"github.com/haorendashu/simdoc/sdcomm"
	"io/ioutil"
	"strings"
)

func Doc(w http.ResponseWriter, r *http.Request, m pongo2.Context) {
	docTpl := pongo2.Must(pongo2.FromFile("views/doc.html"))

	path :=  r.URL.Path
	path = strings.Replace(path, "/doc/", "", 1)
	if path == "" {
		path = "README.md"
	}
	path = filepath.Join(sdcomm.RootPath, path)

	file, err := os.Open(path)
	if err != nil {
		sdcomm.Logger.Println(err)
		if os.IsNotExist(err) {
			Return404(w)
			return
		} else {
			Return500(w)
			return
		}
	}
	defer file.Close()

	// check if is dir
	fileInfo, err := os.Stat(path)
	if err != nil {
		Return500(w)
		return
	}
	if fileInfo.IsDir() {
		listPath(w, m, path)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		sdcomm.Logger.Println(err)
		Return500(w)
		return
	}

	var content string
	var name string
	contentType, isTextFile :=  IsTextFile(data)
	if !isTextFile {
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
		return
	} else {
		// TODO maybe should check code
//		if IsMarkdownFile(contentType) {
			data = blackfriday.MarkdownCommon(data)
//		}
		content = string(data)
	}

	name = filepath.Base(path)
	m["name"] = name
	m["content"] = content

	docTpl.ExecuteWriter(m, w)
}
