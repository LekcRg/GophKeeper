package create

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type NoteModel struct {
	help    *help.CreateVaultNote
	actions *actions.Actions
	log     *zap.Logger
	form    *form.Form
}

var (
	noteNameInput = "name"
	noteTextInput = "text"
)

func NewNote(acts *actions.Actions, log *zap.Logger) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Name",
			IsFocus:     true,
			Name:        noteNameInput,
			Valid:       []validation.Rule{validation.Required},
		}),
	}

	textareas := []components.Textarea{
		components.NewTextarea(components.TextareaOpts{
			Name:        noteTextInput,
			Placeholder: "Enter your note here...",
		}),
	}

	buttons := []components.Button{
		{
			Label: "Create",
			Name:  passwordCreateBtn,
		},
	}

	h := help.NewCreateVaultNote()

	return &NoteModel{
		actions: acts,
		form: form.NewForm(form.Opts{
			Inputs:    inputs,
			Buttons:   buttons,
			Textareas: textareas,
			Up:        h.Keys.Up,
			Down:      h.Keys.Down,
		}),
		help: h,
		log:  log,
	}
}

func (m *NoteModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *NoteModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		name := msg.Values[passwordNameInput]
		data := models.VaultNote{
			Text: msg.Values[noteTextInput],
		}

		res, err := m.actions.CreateVaultItem(context.Background(), name, "note", data)
		if err != nil {
			return msgs.ErrorMsg(err)
		}

		return msgs.CreateVaultSuccess{Item: res}
	}
}

func (m *NoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	default:
		switch typeMsg := msg.(type) {
		case msgs.FormSubmitMsg:
			return m, m.handleSubmit(typeMsg)
		default:
			var newMsg tea.Cmd
			m.form, newMsg = m.form.Update(msg)

			return m, newMsg
		}
	}
}

func (m *NoteModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
