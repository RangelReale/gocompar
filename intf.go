package gocompar

import (
	"io"
	"os"
)

//
// Parser interface, implement this to parse a specific file format
//
type ParserIntf interface {
	Parse(reader io.Reader) ([]*Comment, error)
}

//
// Filter interface, this filters the paths that can be scanned (optional)
//
type FilterIntf interface {
	CanParse(dir string, fileinfo os.FileInfo) bool
}
