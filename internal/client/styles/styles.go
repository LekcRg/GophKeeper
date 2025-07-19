package styles

import "github.com/charmbracelet/lipgloss"

//nolint:gochecknoglobals // styles
var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	CursorStyle  = FocusedStyle
	NoStyle      = lipgloss.NewStyle()
)
