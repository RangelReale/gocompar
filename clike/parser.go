package gcp_clike

import (
	"io"

	"bufio"
	"github.com/RangelReale/gocompar"
	"strings"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(reader io.Reader) ([]*gocompar.Comment, error) {

	scanner := bufio.NewScanner(reader)
	pc := newparserComment()

	for scanner.Scan() {
		//fmt.Printf("@@@ %s\n", scanner.Text())

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
	pc.finish()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return pc.comments, nil

}
