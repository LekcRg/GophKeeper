package components

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type BinaryOpts struct {
	Path string
	ID   int
}

type DetailModel struct {
	help       *help.Auth
	helpBinary *help.BinaryDetail
	actions    *actions.Actions
	error      string
	path       string
	savedPath  string
	fields     []Field
	binaryOpts BinaryOpts
	id         int
}

func NewDetail(field []Field, acts *actions.Actions, binaryOpts BinaryOpts) tea.Model {
	return &DetailModel{
		help:       help.NewAuth(),
		helpBinary: help.NewBinaryDetail(),
		fields:     field,
		path:       binaryOpts.Path,
		id:         binaryOpts.ID,
		actions:    acts,
		binaryOpts: binaryOpts,
	}
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) handleDownload() tea.Cmd {
	return func() tea.Msg {
		path, err := m.actions.DownloadBinary(context.Background(), m.binaryOpts.Path, m.binaryOpts.ID)
		if err != nil {
			m.error = err.Error()
		}

		m.savedPath = path

		return ""
	}
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.binaryOpts.ID == 0 {
		return m, nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if key.Matches(keyMsg, m.helpBinary.Keys.Download) {
		return m, m.handleDownload()
	}

	return m, nil
}

func (m *DetailModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	for i := range m.fields {
		field := &m.fields[i]
		b.WriteString(field.View())
		b.WriteRune('\n')
	}

	if m.binaryOpts.ID > 0 {
		b.WriteRune('\n')

		if m.savedPath != "" {
			b.WriteString("file saved to: " + m.savedPath)
		}

		if m.error != "" {
			b.WriteString(m.error)
		}

		b.WriteRune('\n')
		b.WriteString(m.helpBinary.View())
	} else {
		b.WriteRune('\n')
		b.WriteString(m.help.View())
	}

	return b.String()
}
