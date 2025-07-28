package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type SelectAuth struct {
	help help.Model
	keys *SelectAuthKeyMap
}

func NewSelectAuth() *SelectAuth {
	return &SelectAuth{
		keys: &SelectAuthKeyMap{
			Up: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("↓", "move down"),
			),
			Select: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
		help: help.New(),
	}
}

func (m *SelectAuth) View() string {
	return m.help.View(m.keys)
}

type SelectAuthKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

func (k *SelectAuthKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

func (k *SelectAuthKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
