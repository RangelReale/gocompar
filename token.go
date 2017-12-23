package gocompar

// Comment flags
type Flags int

const (
	// The source comment was multi-line
	MULTI_LINE Flags = 1 << iota

	// The source comment started at the beginning of the line
	START_FULL_LINE

	// The source comment was a collection of single line comments in sequence concatenaed
	CONCAT_MULTI
)

// Parsed comment
type Comment struct {
	Line   int
	Column int
	Text   string
	Flags  Flags
}

// Separates the comments by file
type CommentFile struct {
	Filename string
	Comments []*Comment
}
