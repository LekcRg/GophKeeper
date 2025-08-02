package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type BinaryDetail struct {
	help help.Model
	Keys *BinaryDetailKeyMap
}

func NewBinaryDetail() *BinaryDetail {
	return &BinaryDetail{
		Keys: &BinaryDetailKeyMap{
			Up:       Up,
			Down:     Down,
			Select:   Select,
			Quit:     Quit,
			Download: Download,
		},
		help: help.New(),
	}
}

func (m *BinaryDetail) View() string {
	return m.help.View(m.Keys)
}

type BinaryDetailKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Quit     key.Binding
	Download key.Binding
}

func (k *BinaryDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit, k.Download}
}

func (k *BinaryDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
