package terminal

import (
	"strings"
	"unicode"
)

func FormatTextBreak(text string, width int, tabSize int) []string {
	text = strings.Replace(text, "\r\n", "\n", -1)
	text = strings.Replace(text, "\r", "\n", -1)
	text = strings.Replace(text, "\t", strings.Repeat(" ", tabSize), -1)
	preLines := strings.Split(text, "\n")
	blankLine := strings.Repeat(" ", width)
	lines := make([]string, 0)
	for _, line := range preLines {
		if len(line) == 0 { //newline only
			lines = append(lines, blankLine)
			continue
		}
		for len(line) > 0 {
			// find last space at or before width
			spaceIndex := -1
			spaceCnt := -1
			var cnt, index int
			var r rune
			for index, r = range line {
				if unicode.IsPrint(r) {
					cnt++
					if r == ' ' {
						spaceIndex = index
						spaceCnt = cnt
					}
					if cnt >= width {
						break
					}
				}
			}
			if cnt == 0 { // no printable chars in line, add blankLine to any non-printable
				lines = append(lines, line+blankLine)
				break
			}
			if spaceIndex >= 0 && cnt == width {
				lines = append(lines, line[:spaceIndex]+strings.Repeat(" ", width-spaceCnt+1))
				line = strings.TrimLeft(line[spaceIndex:], " ")
				continue
			}
			// if we get here, then just break at width
			lines = append(lines, line[:index+1]+strings.Repeat(" ", width-cnt))
			line = strings.TrimLeft(line[index+1:], " ")
		}
	}
	return lines
}
