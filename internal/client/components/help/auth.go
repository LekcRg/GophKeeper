package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Auth struct {
	help help.Model
	Keys *AuthKeyMap
}

func NewAuth() *Auth {
	return &Auth{
		Keys: &AuthKeyMap{
			Up:   UpShift,
			Down: DownShift,
			Back: Back,
			Quit: Quit,
		},
		help: help.New(),
	}
}

func (au *Auth) View() string {
	return au.help.View(au.Keys)
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
