package terminal

import (
	"fmt"
	"strings"

	"github.com/buger/goterm"
)

type Area struct {
	terminal Terminal

	content string

	height int
}

// NewArea returns a new Area.
func NewArea() Area {
	return Area{
		terminal: NewDefault(),
		height:   0,
	}
}

// Clear clears the content of the Area.
func (area *Area) Clear() {
	area.terminal.StartOfLine()

	for area.height > 0 {
		area.terminal.ClearLine()
		area.terminal.PrevLine(1)
		area.height--
	}
	area.terminal.ClearLine()
	goterm.Flush()
}

// Update overwrites the content of the Area.
func (area *Area) Update(content string) {
	area.Clear()

	w := goterm.Width()
	lines := FormatTextBreak(strings.TrimSuffix(content, "\n"), w, 3)

	for _, line := range lines {
		fmt.Println(line)
	}

	area.content = content
	area.height = len(lines)
	goterm.Flush()
}
