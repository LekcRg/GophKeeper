package views

import tea "github.com/charmbracelet/bubbletea"

type TokenAuthModel struct {
	tea.Model
}

func NewTokenAuth() *TokenAuthModel {
	return &TokenAuthModel{}
}

func (m *TokenAuthModel) Init() tea.Cmd {
	return nil
}

func (m *TokenAuthModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m *TokenAuthModel) View() string {
	return "Token auth"
}
