package main

import (
	"fmt"
	"os"

	"github.com/LekcRg/GophKeeper/internal/client/views"
	tea "github.com/charmbracelet/bubbletea"
)

type CurrentView int

const (
	auth CurrentView = iota
	list
)

type model struct {
	auth tea.Model
	view CurrentView
}

func initialModel() model {
	m := model{
		auth: views.NewAuth(),
		view: auth,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	switch m.view {
	case auth:
		var cmd tea.Cmd
		m.auth, cmd = m.auth.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.view {
	case auth:
		return m.auth.View()
	}

	return "Error, ctrl+c to quit"
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
