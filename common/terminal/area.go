package terminal

import (
	"fmt"
	"strings"

	"github.com/buger/goterm"
)

type Area struct {
	content string

	height int
}

// NewArea returns a new Area.
func NewArea() Area {
	return Area{
		height: 0,
	}
}

// Clear (lazy) instead of clearing the content of the Area,
// moves the cursor to the beginning of the content to overwrite it.
func (area *Area) Clear() {
}

const CSI = "\x1b\x5b"

// CUU - Cursor Up
func CursorUp(n int) string {
	return fmt.Sprintf("%s%dA", CSI, n)
}

func CursorHorizontalAbsolute(n int) string {
	return fmt.Sprintf("%s%dG", CSI, n)
}

// Update overwrites the content of the Area.
func (area *Area) Update(content string) {

	goterm.Print(CursorHorizontalAbsolute(0))

	if area.height != 0 {
		goterm.Print(CursorUp(area.height))
	}

	w := goterm.Width()
	lines := FormatTextBreak(strings.TrimSuffix(content, "\n"), w, 3)
	//lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")

	for _, line := range lines {
		goterm.Println(line)
	}

	area.content = content
	area.height = len(lines)
	goterm.Flush()
}
