package gcp_clike

import (
	"strings"
	"testing"
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

/**
 * Multi-line comment
 */


`
)

func TestGOParseSingleLine(t *testing.T) {
	p := NewParser()
	c, err := p.Parse(strings.NewReader(gosouce))
	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 3 {
		t.Fatalf("Must have 3 comments extracted, have %d", len(c))
	}

	if strings.TrimSpace(c[0].Text) != "Single-line comment" {
		t.Fatal("Comment[0] text mismatch")
	}

	if !strings.HasPrefix(strings.TrimSpace(c[1].Text), "Multi-line single-line comment") {
		t.Fatalf("Comment[1] text mismatch [%s]", strings.TrimSpace(c[1].Text))
	}

	if !strings.Contains(strings.TrimSpace(c[2].Text), "Multi-line comment") {
		t.Fatalf("Comment[1] text mismatch [%s]", strings.TrimSpace(c[2].Text))
	}

}
