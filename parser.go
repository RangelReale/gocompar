package gocompar

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

//
// Parses comments from source files
//
type Parser struct {
	Parser ParserIntf
	Filter FilterIntf

	Comments []*CommentFile
}

func NewParser(parser ParserIntf, filter FilterIntf) *Parser {
	return &Parser{
		Parser: parser,
		Filter: filter,
	}
}

// Parses an io.Reader
func (p *Parser) Parse(reader io.Reader, filename string) error {

	if p.Parser == nil {
		return errors.New("No parser interface specified")
	}

	c, err := p.Parser.Parse(reader)
	if err != nil {
		return err
	}

	if len(c) > 0 {
		p.Comments = append(p.Comments, &CommentFile{
			Filename: filename,
			Comments: c,
		})
	}

	return nil
}

// Parses a file
func (p *Parser) ParseFile(filename string) error {

	if p.Parser == nil {
		return errors.New("No parser interface specified")
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.Parse(f, filename)
}

// Parses a path, using FilterIntf to select which ones
func (p *Parser) ParseDir(dir string) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if p.Filter != nil && !p.Filter.CanParse(dir, f) {
			continue
		}

		if !f.IsDir() {
			err = p.ParseFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return err
			}
		} else {
			err = p.ParseDir(filepath.Join(dir, f.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
