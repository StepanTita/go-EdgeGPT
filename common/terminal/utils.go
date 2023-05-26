package terminal

import (
	"fmt"
	"io"
	"os"

	"github.com/buger/goterm"
)

const CSI = "\x1b\x5b"

type Terminal struct {
	target io.Writer
}

func NewDefault() Terminal {
	return Terminal{
		target: os.Stdout,
	}
}

func (t Terminal) NextLine(n int) {
	goterm.MoveCursorDown(n)
}

func (t Terminal) PrevLine(n int) {
	goterm.MoveCursorUp(n)
}

func (t Terminal) ClearLine() {
	fmt.Fprintf(t.target, "%s0K", CSI)
}
