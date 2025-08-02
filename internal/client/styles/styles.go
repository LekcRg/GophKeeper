package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

//nolint:gochecknoglobals // styles
var (
	Gray           = lipgloss.Color("240")
	Green          = lipgloss.Color("#00D084")
	FocusColor     = lipgloss.Color("205")
	FocusColorText = lipgloss.Color("#000")
	ErrorColor     = lipgloss.Color("9")
	FocusedStyle   = lipgloss.NewStyle().Foreground(FocusColor)
	BlurredStyle   = lipgloss.NewStyle().Foreground(Gray)
	ErrorStyle     = lipgloss.NewStyle().Foreground(ErrorColor)
	CursorStyle    = FocusedStyle
	NoStyle        = lipgloss.NewStyle()
	Border         = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Gray)
	FieldLabel = lipgloss.NewStyle().Bold(true).Foreground(FocusColor)
)

var (
	SuccessTitle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)
	TokenBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(1, 0).
			BorderForeground(Green).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#1A1A1A"))
	Instruction = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)
)

func GetTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Gray).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(FocusColorText)).
		Background(FocusColor).
		Bold(false)

	return s
}

func GetTextareaStyles() (textarea.Style, textarea.Style) {
	focused := textarea.Style{
		Base:             NoStyle,
		CursorLine:       FocusedStyle,
		CursorLineNumber: FocusedStyle,
		EndOfBuffer:      lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "254", Dark: "0"}),
		LineNumber:       NoStyle,
		Placeholder:      BlurredStyle,
		Prompt:           NoStyle,
		Text:             NoStyle,
	}
	blurred := textarea.Style{
		Base:             BlurredStyle,
		CursorLine:       BlurredStyle,
		CursorLineNumber: BlurredStyle,
		EndOfBuffer:      NoStyle,
		LineNumber:       BlurredStyle,
		Placeholder:      BlurredStyle,
		Prompt:           NoStyle,
		Text:             BlurredStyle,
	}

	return focused, blurred
}
