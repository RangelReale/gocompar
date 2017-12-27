package gcp_clike

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/RangelReale/gocompar"
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
	is_data := false
	for !is_eof {
		t := l.nextToken()
		var c *gocompar.Comment
		switch t.typ {
		case tokenCComment:
			f := gocompar.MULTI_LINE
			if !t.LineHasData {
				f |= gocompar.START_FULL_LINE
			}
			c = &gocompar.Comment{
				Line:   t.Line,
				Column: t.Column,
				Text:   strings.TrimSpace(t.String()),
				Flags:  f,
			}
		case tokenCPPComment:
			var f gocompar.Flags
			if !t.LineHasData {
				f |= gocompar.START_FULL_LINE
			}
			c = &gocompar.Comment{
				Line:   t.Line,
				Column: t.Column,
				Text:   strings.TrimSpace(t.String()),
				Flags:  f,
			}
		case tokenEOF:
			is_eof = true
			break
		case tokenError:
			return nil, t
		default:
			if strings.TrimSpace(t.String()) != "" {
				is_data = true
			}
		}

		if c != nil {
			p.cleanComment(c)

			// if current and previous comment are single line and started in a full line,
			// and had no data in between, append to previous
			if len(comments) > 0 &&
				!is_data &&
				c.Flags&gocompar.START_FULL_LINE == gocompar.START_FULL_LINE &&
				comments[len(comments)-1].Flags&gocompar.START_FULL_LINE == gocompar.START_FULL_LINE &&
				c.Flags&gocompar.MULTI_LINE != gocompar.MULTI_LINE &&
				comments[len(comments)-1].Flags&gocompar.MULTI_LINE != gocompar.MULTI_LINE {
				if comments[len(comments)-1].Text != "" {
					comments[len(comments)-1].Text += "\n"
				}
				comments[len(comments)-1].Text += c.Text
				comments[len(comments)-1].Flags |= gocompar.CONCAT_MULTI
			} else {
				comments = append(comments, c)
			}
			is_data = false
		}
	}

	return comments, nil
}

// remove spaces, slashes and asterisks from the beginning of the comment
func (p *Parser2) cleanComment(c *gocompar.Comment) {
	var ret bytes.Buffer

	scan := bufio.NewScanner(strings.NewReader(c.Text))

	is_first := true
	find := map[rune]bool{' ': true, '/': true}
	if c.Flags&gocompar.MULTI_LINE == gocompar.MULTI_LINE {
		find['*'] = true
	}
	for scan.Scan() {
		// always trim right spaces
		nline := strings.TrimRightFunc(scan.Text(), unicode.IsSpace)

		// check if there are spaces, slashes and asterisks at the beginning, and remove them
		line := nline
		for li, lr := range nline {
			_, isfind := find[lr]
			//if lr != ' ' && lr != '*' && lr != '/' {
			if !isfind {
				break
			}
			if lr == '*' || lr == '/' {
				line = nline[li+1:]
			}
		}

		// if there is one space at left, remove it
		if len(line) > 0 && line[0] == ' ' {
			line = line[1:]
		}

		//if !is_first || strings.TrimSpace(line) != "" {
		if strings.TrimSpace(line) != "" {
			if !is_first {
				ret.WriteString("\n")
			}
			ret.WriteString(line)
			is_first = false
		}
	}

	c.Text = ret.String()
}
