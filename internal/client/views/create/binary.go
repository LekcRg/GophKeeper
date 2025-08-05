package create

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/errs"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type BinaryModel struct {
	help    *help.Auth
	actions *actions.Actions
	log     *zap.Logger
	form    *form.Form
	path    string
}

const (
	binaryNameInput     = "name"
	choiceFileNameBtn   = "choice"
	createNameBtn       = "create"
	fileNotSelectedText = "File not selected"
)

func NewBinary(acts *actions.Actions, log *zap.Logger) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Name",
			Name:        binaryNameInput,
			IsFocus:     true,
			Valid:       []validation.Rule{validation.Required},
		}),
	}

	buttons := []components.Button{
		{
			Label:          "Choice file",
			Name:           choiceFileNameBtn,
			AdditionalText: fileNotSelectedText,
		},
		{
			Label: "Create",
			Name:  createNameBtn,
		},
	}

	h := help.NewAuth()

	return &BinaryModel{
		actions: acts,
		form: form.NewForm(form.Opts{
			Inputs:  inputs,
			Buttons: buttons,
			Up:      h.Keys.Up,
			Down:    h.Keys.Down,
		}),
		help: h,
		log:  log,
	}
}

func (m *BinaryModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *BinaryModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		if msg.ButtonName == choiceFileNameBtn {
			return msgs.OpenFilePicker{}
		}

		if m.path == "" {
			m.form.HandleError(errs.ErrFileEmpty, "")
			return ""
		}

		res, err := m.actions.CreateBinaryVault(
			context.Background(),
			msg.Values[binaryNameInput],
			m.path,
		)
		if err != nil {
			m.log.Error("CreateBinaryVault error", zap.Error(err))
			return msgs.ErrorMsg(err)
		}

		return msgs.CreateVaultSuccess{Item: res}
	}
}

func (m *BinaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.handleSubmit(typeMsg)
	case msgs.FilepickerSelected:
		m.path = string(typeMsg)
		m.form.ChangeButtonAdditionalText(0, m.path)
	}

	var newCmd tea.Cmd
	m.form, newCmd = m.form.Update(msg)

	return m, newCmd
}

func (m *BinaryModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())

	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
