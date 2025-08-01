package views

import (
	"context"
	"strconv"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/state"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type ListModels struct {
	help   *help.List
	log    *zap.Logger
	state  *state.State
	error  string
	table  table.Model
	loaded bool
}

func NewList(l *zap.Logger, st *state.State) *ListModels {
	var (
		ColumnIDWitdh      = 0
		ColumnNameWidth    = 20
		ColumnTypeWidth    = 10
		ColumnUpdatedWidth = 20
		RowHeight          = 7
	)

	columns := []table.Column{
		{Title: "ID", Width: ColumnIDWitdh},
		{Title: "Name", Width: ColumnNameWidth},
		{Title: "Type", Width: ColumnTypeWidth},
		{Title: "Updated at", Width: ColumnUpdatedWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(RowHeight),
	)

	t.SetStyles(styles.GetTableStyles())

	return &ListModels{
		table: t,
		log:   l,
		help:  help.NewList(),
		state: st,
	}
}

func (m *ListModels) load() tea.Msg {
	err := m.state.LoadVault(context.Background())
	if err != nil {
		m.log.Error("state LoadVault error", zap.Error(err))
		return msgs.ErrorMsg(err)
	}

	return msgs.ListLoaded{}
}

func (m *ListModels) Init() tea.Cmd {
	if m.loaded {
		return tea.WindowSize()
	}

	return tea.Batch(m.load, tea.WindowSize())
}

func (m *ListModels) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		margin := 6
		m.table.SetHeight(msg.Height - margin)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.help.Keys.Select):
			if len(m.state.Table) > 0 {
				selectedID, err := strconv.Atoi(m.table.SelectedRow()[0])
				if err != nil {
					errText := "Selected vault ID is not int"
					m.log.Error(errText, zap.Error(err))
					m.error = errText

					return m, nil
				}

				return m, func() tea.Msg {
					return msgs.SelectVaultItem(selectedID)
				}
			}
		case key.Matches(msg, m.help.Keys.Create):
			return m, func() tea.Msg {
				return msgs.ToCreateVaultItem{}
			}
		}
	case msgs.ListLoaded:
		m.table.SetRows(m.state.Table)
	case msgs.ErrorMsg:
		m.error = msg.Error()
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *ListModels) View() string {
	var b strings.Builder

	b.WriteString(styles.Border.Render(m.table.View()))
	b.WriteRune('\n')
	b.WriteString(styles.ErrorStyle.Render(m.error))
	b.WriteString("\n\n")
	b.WriteString(m.help.View())
	b.WriteString("\n")

	return b.String()
}
