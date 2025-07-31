package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type CreateVaultNote struct {
	help help.Model
	Keys *CreateVaultNoteKeyMap
}

func NewCreateVaultNote() *CreateVaultNote {
	return &CreateVaultNote{
		Keys: &CreateVaultNoteKeyMap{
			Up:     UpShiftOnly,
			Down:   DownShiftOnly,
			Select: Select,
			Quit:   Quit,
		},
		help: help.New(),
	}
}

func (m *CreateVaultNote) View() string {
	return m.help.View(m.Keys)
}

type CreateVaultNoteKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

func (k *CreateVaultNoteKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

func (k *CreateVaultNoteKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
