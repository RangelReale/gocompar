package gcp_clike

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"

	"github.com/RangelReale/gocompar"
)

//
// State machine to extract comments from golang file
// TODO: rework this mess
//
type parserComment struct {
	comments              []*gocompar.Comment
	is_comment            bool
	is_comment_singleline bool
	is_slash              bool
	is_asterisk           bool
	curcomment            bytes.Buffer
	comment_line_ct       int
	comment_column_ct     int
	comment_line_had_data int

	line_has_data    int
	noncomment_lines int
	line_ct          int
	column_ct        int
}

func newparserComment() *parserComment {
	return &parserComment{
		comments: make([]*gocompar.Comment, 0),
	}
}

func (pc *parserComment) addRune(c rune) {

	pc.column_ct++

	current_is_comment := false
	current_end_comment := false
	if c == '/' {
		if pc.is_slash && !pc.is_comment {
			current_is_comment = true
			pc.is_comment_singleline = true
		}
		if pc.is_asterisk && pc.is_comment {
			current_end_comment = true
		}
		pc.is_slash = true
		pc.is_asterisk = false
	} else if c == '*' {
		if pc.is_slash && !pc.is_comment {
			current_is_comment = true
			pc.is_comment_singleline = false
		}
		pc.is_slash = false
		pc.is_asterisk = true
	} else {
		pc.is_slash = false
		pc.is_asterisk = false
	}

	if current_is_comment {
		pc.is_comment = true
		pc.line_has_data--
		pc.noncomment_lines--
		pc.comment_line_had_data = pc.line_has_data
		pc.comment_line_ct = pc.line_ct
		pc.comment_column_ct = pc.column_ct - 2
	} else if current_end_comment {
		pc.addComment()
		pc.is_slash = false
		pc.is_asterisk = false
	} else if pc.is_comment {
		pc.curcomment.WriteRune(c)
	} else {
		if !unicode.IsSpace(c) {
			//fmt.Printf("NOT IS SPACE: %v [%x]\n", c, c)
			pc.line_has_data++
			pc.noncomment_lines++
		}
	}

}

func (pc *parserComment) addNewLine() {

	if pc.is_comment {
		if pc.is_comment_singleline {
			pc.addComment()
		} else {
			pc.curcomment.WriteRune('\n')
		}
	}
	pc.column_ct = 0
	pc.line_ct++
	pc.line_has_data = 0
	pc.is_slash = false
	pc.is_asterisk = false

}

func (pc *parserComment) addComment() {
	// concatenate sequential single line comments
	if len(pc.comments) > 0 && pc.noncomment_lines == 0 &&
		(pc.comments[len(pc.comments)-1].Flags&gocompar.START_FULL_LINE) == gocompar.START_FULL_LINE &&
		(pc.comments[len(pc.comments)-1].Flags&gocompar.MULTI_LINE) != gocompar.MULTI_LINE {
		// add comment to previous line
		pc.comments[len(pc.comments)-1].Text += "\n" + pc.cleanedCurComment()
		pc.comments[len(pc.comments)-1].Flags |= gocompar.CONCAT_MULTI
	} else {
		var fl gocompar.Flags
		if !pc.is_comment_singleline {
			fl |= gocompar.MULTI_LINE
		}
		if pc.comment_line_had_data < 2 {
			fl |= gocompar.START_FULL_LINE
		}
		pc.comments = append(pc.comments, &gocompar.Comment{
			Line:   pc.comment_line_ct + 1,
			Column: pc.comment_column_ct,
			Text:   pc.cleanedCurComment(),
			Flags:  fl,
		})
	}
	pc.curcomment.Reset()
	pc.is_comment = false
	pc.noncomment_lines = 0
}

func (pc *parserComment) finish() {
	if pc.is_comment {
		pc.addComment()
	}
}

func (pc *parserComment) cleanedCurComment() string {

	var ret bytes.Buffer

	scan := bufio.NewScanner(bytes.NewReader(pc.curcomment.Bytes()))

	is_first := true
	for scan.Scan() {

		line := strings.TrimSpace(scan.Text())

		if !pc.is_comment_singleline {
			// remove asterisks from beginning
			line = strings.TrimLeftFunc(line, func(c rune) bool {
				if unicode.IsSpace(c) || c == '*' {
					return true
				}
				return false
			})
		}

		if !is_first || line != "" {
			if !is_first {
				ret.WriteString("\n")
			}
			ret.WriteString(line)
			is_first = false
		}
	}

	return strings.TrimSpace(ret.String())
}
