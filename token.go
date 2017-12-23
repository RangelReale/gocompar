package gocompar

type Token int

const (
	COMMENT Token = iota
	EOF
)

type Flags int

const (
	MULTI_LINE Flags = 1 << iota
	START_FULL_LINE
	CONCAT_MULTI
)

type Comment struct {
	Line   int
	Column int
	Text   string
	Flags  Flags
}
