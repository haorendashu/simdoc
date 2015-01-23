package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	rootPath string
)

func main() {
	// rootPath
	file, _ := exec.LookPath(os.Args[0])
	rootPath, _ = filepath.Abs(file)
	rootPath, _ = filepath.Split(rootPath)
	rootPath = filepath.Join(rootPath, "docs")

}
