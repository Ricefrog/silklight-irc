package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles with just a color.
var MagentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
var RedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

var Separator = RedStyle.Render("╬")

var BorderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

// Styles for nicknames
var NickStyleSelf = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

var nickCols = []string{"2", "4", "5", "6", "10", "12", "13", "14"}

func NickToStyle(nick string) lipgloss.Style {
	sum := 0
	for _, ch := range nick {
		sum += int(ch)
	}
	sum %= len(nickCols)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(nickCols[sum]))
}
