package gcp_clike

import (
	"strings"
	"testing"

	"github.com/RangelReale/gocompar"
)

const (
	gosouce = `

// Single-line comment

var (
	x = "test value"
)

// Multi-line single-line comment
// in following lines

func main() {

}

/*************
 * Multi-line comment
 ******/

var a string // Comment on a variable

`
)

func TestGOParseSingleLine(t *testing.T) {
	p := NewParser()
	c, err := p.Parse(strings.NewReader(gosouce))
	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 4 {
		t.Fatalf("Must have 3 comments extracted, have %d", len(c))
	}

	if c[0].Text != "Single-line comment" {
		t.Fatal("Comment[0] text mismatch")
	}

	if c[0].Flags&gocompar.MULTI_LINE == gocompar.MULTI_LINE {
		t.Fatalf("Comment[0] flag MULTI_LINE should NOT had been set [%s]", c[0].Text)
	}

	if c[1].Text != "Multi-line single-line comment\nin following lines" {
		t.Fatalf("Comment[1] text mismatch [%s]", c[1].Text)
	}
	if c[1].Flags&gocompar.CONCAT_MULTI != gocompar.CONCAT_MULTI {
		t.Fatalf("Comment[1] flag CONCAT_MULTI should been set [%s]", c[1].Text)
	}

	if c[2].Text != "Multi-line comment" {
		t.Fatalf("Comment[2] text mismatch [%s]", c[2].Text)
	}

	if c[2].Flags&gocompar.MULTI_LINE != gocompar.MULTI_LINE {
		t.Fatalf("Comment[2] flag MULTI_LINE should been set [%s]", c[2].Text)
	}

	if c[3].Text != "Comment on a variable" {
		t.Fatalf("Comment[3] text mismatch [%s]", c[3].Text)
	}

	if c[3].Flags&gocompar.START_FULL_LINE == gocompar.START_FULL_LINE {
		t.Fatalf("Comment[3] flag START_FULL_LINE should NOT had been set [%s]", c[3].Text)
	}

}
