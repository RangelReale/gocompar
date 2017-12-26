package gcp_clike

import (
	"io"

	"fmt"
	"github.com/RangelReale/gocompar"
	"io/ioutil"
)

//
// Golang comments extractor
//
type Parser2 struct {
}

func NewParser2() *Parser2 {
	return &Parser2{}
}

func (p *Parser2) Parse(reader io.Reader) ([]*gocompar.Comment, error) {

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	comments := make([]*gocompar.Comment, 0)

	l := lex(data)
	is_eof := false
	for !is_eof {
		t := l.nextToken()
		var c *gocompar.Comment
		switch t.typ {
		case tokenCComment:
			c = &gocompar.Comment{
				Line:   0,
				Column: 0,
				Text:   t.String(),
				Flags:  gocompar.MULTI_LINE,
			}
		case tokenCPPComment:
			c = &gocompar.Comment{
				Line:   0,
				Column: 0,
				Text:   t.String(),
				Flags:  0,
			}
		case tokenEOF:
			is_eof = true
			break
		case tokenError:
			return nil, t
		default:

		}

		if c != nil {
			comments = append(comments, c)
		}
	}

	return comments, nil
}

func (p *Parser2) Dump(reader io.Reader) ([]*gocompar.Comment, error) {

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	l := lex(data)
	for {
		t := l.nextToken()
		switch t.typ {
		case tokenCComment:
			fmt.Printf("@@ TOKEN: C COMMENT\n")
		case tokenCPPComment:
			fmt.Printf("@@ TOKEN: CPP COMMENT\n")
		case tokenQuotedText:
			fmt.Printf("@@ TOKEN: QUOTED TEXT\n")
		case tokenEOF:
			fmt.Printf("@@ TOKEN: EOF\n")
			break
		case tokenError:
			return nil, t
		default:
			fmt.Printf("@@ TOKEN: OTHER\n")
		}
		fmt.Printf("----------- TOKEN DATA BEGIN -----------\n")
		fmt.Printf(t.String())
		fmt.Printf("\n----------- TOKEN DATA END -----------\n")
	}

	return nil, nil
}
