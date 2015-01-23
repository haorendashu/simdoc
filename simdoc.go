package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
	"flag"
	"net/http"
	"log"
	"github.com/flosch/pongo2"
	"runtime/debug"
	"time"
	"github.com/haorendashu/simdoc/routers"
	"github.com/haorendashu/simdoc/sdcomm"
)

var (
	rootPath string
	port string
)

func main() {
	// rootPath
	file, _ := exec.LookPath(os.Args[0])
	rootPath, _ = filepath.Abs(file)
	rootPath, _ = filepath.Split(rootPath)
	rootPath = filepath.Join(rootPath, "docs")
	flag.StringVar(&rootPath, "rootPath", rootPath, "docs root path.")
	flag.StringVar(&port, "port", "8080", "server port.")
	sdcomm.RootPath = rootPath

	sdcomm.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", handlerWrapper(routers.Index))
	r.HandleFunc("/list/", handlerWrapper(routers.List))
	r.PathPrefix("/doc/").HandlerFunc(handlerWrapper(routers.Doc))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	serveSingle("/favicon.ico", "./static/favicon.ico")
	http.Handle("/", r)

	sdcomm.Logger.Printf("server would start at port %s.\n", port)
	sdcomm.Logger.Printf("doc base path is %s.\n", rootPath)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		sdcomm.Logger.Fatalln(err)
	}
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filename)
		})
}

// HTTP Handler 包装，完成共性处理.
//
// 共性处理：
//
//  1. panic recover
//  2. 请求计时
func handlerWrapper(f func(w http.ResponseWriter, r *http.Request, m pongo2.Context)) func(w http.ResponseWriter, r *http.Request) {
	handler := commonWrapper(f)
	handler = panicRecover(handler)
	handler = stopwatch(handler)

	return handler
}

// Handler 包装请求计时.
func stopwatch(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			uas := r.Header["User-Agent"]
			userAgent := ""
			if len(uas) > 0 {
				userAgent = uas[0]
			}
			sdcomm.Logger.Printf("[%s] [%s] [%s] [%s]", r.RequestURI, r.Method, time.Since(start), userAgent)
		}()

		// Handler 处理
		handler(w, r)
	}
}

// Handler 包装 recover panic.
func panicRecover(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// panic 恢复.
		defer func() {
			if re := recover(); re != nil {
				sdcomm.Logger.Printf("PANIC RECOVERED:\n %v, %s", re, debug.Stack())
			}
		}()

		// Handler 处理
		handler(w, r)
	}
}

func commonWrapper(f func(w http.ResponseWriter, r *http.Request, m pongo2.Context)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		m := pongo2.Context{}
		f(w, r, m)
	}
}
