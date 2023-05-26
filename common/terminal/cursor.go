package terminal

import "fmt"

// StartOfLine moves the cursor to the start of the current line.
func (t Terminal) StartOfLine() {
	t.HorizontalAbsolute(0)
}

// StartOfLineDown moves the cursor down by n lines, then moves to cursor to the start of the line.
func (t Terminal) StartOfLineDown(n int) {
	t.NextLine(n)
	t.StartOfLine()
}

// StartOfLineUp moves the cursor up by n lines, then moves to cursor to the start of the line.
func (t Terminal) StartOfLineUp(n int) {
	t.PrevLine(n)
	t.StartOfLine()
}

// HorizontalAbsolute moves the cursor to n horizontally.
// The position n is absolute to the start of the line.
func (t Terminal) HorizontalAbsolute(n int) {
	n += 1 // Moves the line to the character after n
	fmt.Fprintf(t.target, "%s%dG", CSI, n)
}
