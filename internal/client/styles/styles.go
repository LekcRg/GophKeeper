package styles

import "github.com/charmbracelet/lipgloss"

//nolint:gochecknoglobals // styles
var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	CursorStyle  = FocusedStyle
	NoStyle      = lipgloss.NewStyle()
	Green        = lipgloss.Color("#00D084")
	Gray         = lipgloss.Color("#888888")
)

//nolint:mnd // styles
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
