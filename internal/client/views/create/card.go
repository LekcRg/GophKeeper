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
	"go.uber.org/zap"
)

type CardModel struct {
	help    *help.Auth
	actions *actions.Actions
	log     *zap.Logger
	form    *form.Form
}

const (
	cardNameInput   = "name"
	cardNumberInput = "number"
	cardExpireInput = "expire"
	cardCVVInput    = "cvv"
)

func NewCard(acts *actions.Actions, log *zap.Logger) tea.Model {
	const (
		cardMaxLen = 19
		expMaxLen  = 5
		cvvMaxLen  = 3
	)

	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Name",
			Name:        cardNameInput,
			IsFocus:     true,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Number",
			Name:        cardNumberInput,
			CharLimit:   cardMaxLen,
			Type:        components.InputTypeCardNumber,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Expire",
			Name:        cardExpireInput,
			CharLimit:   expMaxLen,
			Type:        components.InputTypeCardExpire,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "CVV",
			Name:        cardCVVInput,
			CharLimit:   cvvMaxLen,
			Type:        components.InputTypeCardNumber,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Create",
			Name:  passwordCreateBtn,
		},
	}

	h := help.NewAuth()

	return &CardModel{
		actions: acts,
		form: form.NewForm(form.FormOpts{
			Inputs:  inputs,
			Buttons: buttons,
			Up:      h.Keys.Up,
			Down:    h.Keys.Down,
		}),
		help: h,
		log:  log,
	}
}

func (m *CardModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CardModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		name := msg.Values[cardNameInput]
		data := models.VaultItemDataCard{
			Number: msg.Values[cardNumberInput],
			Exp:    msg.Values[cardExpireInput],
			CVV:    msg.Values[cardCVVInput],
		}

		res, err := m.actions.CreateVaultItem(context.Background(), name, "card", data)
		if err != nil {
			return msgs.ErrorMsg(err)
		}

		return msgs.CreateVaultSuccess{Item: res}
	}
}

func (m *CardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *CardModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
