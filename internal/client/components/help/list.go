package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type List struct {
	help help.Model
	Keys *ListKeyMap
}

func NewList() *List {
	return &List{
		Keys: &ListKeyMap{
			Up:     UpK,
			Down:   DownJ,
			Select: Select,
			Quit:   Quit,
		},
		help: help.New(),
	}
}

func (m *List) View() string {
	return m.help.View(m.Keys)
}

type ListKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

func (k *ListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

func (k *ListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
