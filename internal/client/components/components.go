package components

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	View() string
	Update() (Component, tea.Cmd)
	Focus() Component
	Blur() Component
}
