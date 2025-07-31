package help

import "github.com/charmbracelet/bubbles/key"

var (
	Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	)
	Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	)
	UpShift = key.NewBinding(
		key.WithKeys("up", "shift+tab"),
		key.WithHelp("↑/Shift+Tab", "move up"),
	)
	DownShift = key.NewBinding(
		key.WithKeys("down", "tab"),
		key.WithHelp("↓/Tab", "move down"),
	)
	UpK = key.NewBinding(
		key.WithKeys("up", "k", "K"),
		key.WithHelp("↑/k", "move up"),
	)
	DownJ = key.NewBinding(
		key.WithKeys("down", "j", "J"),
		key.WithHelp("↓/j", "move down"),
	)
	Back = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "back"),
	)
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	Logout = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "logout"),
	)
	Select = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)
	Create = key.NewBinding(
		key.WithKeys("c", "C"),
		key.WithHelp("c", "create"),
	)
)
