package terminal

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	appName      lipgloss.Style
	cliArgs      lipgloss.Style
	comment      lipgloss.Style
	cyclingChars lipgloss.Style
	ErrorHeader  lipgloss.Style
	ErrorDetails lipgloss.Style
	flag         lipgloss.Style
	flagComma    lipgloss.Style
	flagDesc     lipgloss.Style
	inlineCode   lipgloss.Style
	link         lipgloss.Style
	pipe         lipgloss.Style
	quote        lipgloss.Style
	Prefix       lipgloss.Style
}

func NewStyles(r *lipgloss.Renderer) (s Styles) {
	s.appName = r.NewStyle().Bold(true)
	s.cliArgs = r.NewStyle().Foreground(lipgloss.Color("#585858"))
	s.comment = r.NewStyle().Foreground(lipgloss.Color("#757575"))
	s.cyclingChars = r.NewStyle().Foreground(lipgloss.Color("#FF87D7"))
	s.ErrorHeader = r.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#FF5F87")).Bold(true).Padding(0, 1).SetString("ERROR")
	s.ErrorDetails = s.comment.Copy()
	s.flag = r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#00B594", Dark: "#3EEFCF"}).Bold(true)
	s.flagComma = r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#5DD6C0", Dark: "#427C72"}).SetString(",")
	s.flagDesc = s.comment.Copy()
	s.inlineCode = r.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Background(lipgloss.Color("#3A3A3A")).Padding(0, 1)
	s.link = r.NewStyle().Foreground(lipgloss.Color("#00AF87")).Underline(true)
	s.quote = r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#FF71D0", Dark: "#FF78D2"})
	s.pipe = r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#8470FF", Dark: "#745CFF"})
	s.Prefix = r.NewStyle().Bold(true).Foreground(lipgloss.Color("#32CD32"))
	return s
}
