package create

import (
	"os"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type filePicker struct {
	filepicker filepicker.Model
	err        error
}

func NewFilePicker() tea.Model {
	fp := filepicker.New()

	var err error
	fp.CurrentDirectory, err = os.UserHomeDir()
	fp.FileAllowed = true
	fp.AutoHeight = false

	m := filePicker{
		filepicker: fp,
		err:        err,
	}

	return &m
}

func (m *filePicker) Init() tea.Cmd {
	var err error

	m.filepicker.CurrentDirectory, err = os.UserHomeDir()
	if err != nil {
		m.err = err
	} else {
		m.err = nil
	}

	return tea.Batch(m.filepicker.Init(), tea.WindowSize())
}

func (m *filePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case tea.WindowSizeMsg:
		marginBottom := 6

		m.filepicker.SetHeight(typedMsg.Height - marginBottom)
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		return m, func() tea.Msg {
			return msgs.FilepickerSelected(path)
		}
	}

	return m, cmd
}

func (m *filePicker) View() string {
	var s strings.Builder

	s.WriteString("\n  " + m.filepicker.CurrentDirectory)
	s.WriteString("\n\n" + m.filepicker.View())

	errText := ""
	if m.err != nil {
		errText = m.err.Error()
	}

	s.WriteString(m.filepicker.Styles.DisabledFile.Render(errText))

	return s.String()
}
