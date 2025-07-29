package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Register struct {
	help help.Model
	keys *AuthKeyMap
}

func NewRegister() *Register {
	return &Register{
		keys: &AuthKeyMap{
			Up: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "tab"),
				key.WithHelp("↓", "move down"),
			),
			Back: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("Esc", "back"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
		help: help.New(),
	}
}

func (au *Register) View() string {
	return au.help.View(au.keys)
}

type AuthKeyMap struct {
	Up   key.Binding
	Down key.Binding
	Back key.Binding
	Quit key.Binding
}

func (k *AuthKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back, k.Quit}
}

func (k *AuthKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
