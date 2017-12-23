package gcp_clike

import (
	"io"

	"bufio"
	"github.com/RangelReale/gocompar"
	"strings"
)

//
// Golang comments extractor
//
type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(reader io.Reader) ([]*gocompar.Comment, error) {

	scanner := bufio.NewScanner(reader)
	pc := newparserComment()

	for scanner.Scan() {
		r := bufio.NewReader(strings.NewReader(scanner.Text()))
		for {
			c, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, err
				}
			}
			pc.addRune(c)
		}
		pc.addNewLine()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	pc.finish()

	return pc.comments, nil

}
