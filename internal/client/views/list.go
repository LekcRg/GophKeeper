package views

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ListModels struct {
	table table.Model
}

func NewList() *ListModels {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 10},
		{Title: "Created at", Width: 20},
		{Title: "Updated at", Width: 20},
	}

	rows := []table.Row{
		{"Google", "Password", "21 January 2023", "3 September 2023"},
		{"Work Email", "Password", "12 February 2022", "12 February 2022"},
		{"Bank Card", "Card", "8 March 2024", "10 April 2024"},
		{"GitHub", "Password", "17 July 2021", "19 November 2022"},
		{"Server SSH Key", "Note", "25 May 2023", "25 May 2023"},
		{"Wi-Fi Setup", "Note", "9 October 2022", "15 January 2023"},
		{"Dropbox", "Password", "3 December 2020", "3 December 2020"},
		{"Driver License", "Binary", "7 April 2023", "7 April 2023"},
		{"Netflix", "Password", "16 June 2022", "18 August 2023"},
		{"Tax Documents", "Binary", "28 February 2024", "28 February 2024"},
		{"Company VPN", "Password", "11 November 2021", "4 May 2023"},
		{"Medical Card", "Card", "14 July 2023", "14 July 2023"},
		{"Steam", "Password", "6 August 2020", "12 October 2022"},
		{"Server Backup Key", "Binary", "30 September 2022", "30 September 2022"},
		{"LinkedIn", "Password", "19 March 2021", "20 March 2023"},
		{"GPG Private Key", "Note", "2 May 2022", "2 May 2022"},
		{"Crypto Wallet Seed", "Note", "13 December 2023", "13 December 2023"},
		{"Amazon", "Password", "27 January 2022", "2 February 2024"},
		{"Insurance Card", "Card", "22 June 2024", "22 June 2024"},
		{"Work Contracts", "Binary", "1 March 2023", "10 March 2023"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &ListModels{t}
}

func (m *ListModels) Init() tea.Cmd {
	return nil
}

func (m *ListModels) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *ListModels) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
