package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type AuthHelp struct {
	keys AuthKeyMap
	help help.Model
}

func NewAuthHelp() AuthHelp {
	return AuthHelp{
		keys: AuthKeyMap{
			Up: key.NewBinding(
				key.WithKeys("up", "shift+tab"),
				key.WithHelp("↑/Shift+Tab", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "tab"),
				key.WithHelp("↓/Tab", "move down"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
		help: help.New(),
	}
}

func (au AuthHelp) View() string {
	return au.help.View(au.keys)
}

type AuthKeyMap struct {
	Up   key.Binding
	Down key.Binding
	Quit key.Binding
}

func (k AuthKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k AuthKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
