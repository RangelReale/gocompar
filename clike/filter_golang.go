package gcp_clike

import (
	"os"
	"path/filepath"
	"strings"
)

//
// Filter for golang source files, excluding common subdirectories
//
type Filter_Golang struct {
}

func (p *Filter_Golang) CanParse(dir string, fileinfo os.FileInfo) bool {

	if !fileinfo.IsDir() && strings.ToLower(filepath.Ext(fileinfo.Name())) != ".go" {
		return false
	} else if fileinfo.IsDir() && (fileinfo.Name() == "pkg" || fileinfo.Name() == "vendor" || strings.HasPrefix(fileinfo.Name(), ".")) {
		return false
	}

	return true
}
