package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type SelectAuth struct {
	help help.Model
	Keys *SelectAuthKeyMap
}

func NewSelectAuth() *SelectAuth {
	return &SelectAuth{
		Keys: &SelectAuthKeyMap{
			Up:     Up,
			Down:   Down,
			Select: Select,
			Quit:   Quit,
		},
		help: help.New(),
	}
}

func (m *SelectAuth) View() string {
	return m.help.View(m.Keys)
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
