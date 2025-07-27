package views

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type registerViewErrors struct {
	Login          string
	Password       string
	CryptoPassword string
	Message        string
}

type RegisterModel struct {
	tea.Model
	help          *help.Register
	actions       *actions.Actions
	log           *zap.Logger
	ResponseError models.UserReq
	Response      models.APIKeyRes
	errors        registerViewErrors
	inputs        []components.TextInput
	buttons       []components.Button
	focusIndex    int
}

func NewAuth(acts *actions.Actions, log *zap.Logger) *RegisterModel {
	m := RegisterModel{
		inputs: []components.TextInput{
			components.NewTextInput(components.TextInputOpts{
				Placeholder: "Login",
				IsFocus:     true,
				Name:        "login",
			}),
			components.NewTextInput(components.TextInputOpts{
				Placeholder: "Auth password",
				IsPassword:  true,
				Name:        "password",
			}),
			components.NewTextInput(components.TextInputOpts{
				Placeholder: "Enctyption password",
				IsPassword:  true,
				Name:        "crypto-password",
			}),
		},
		buttons: []components.Button{
			{
				Label: "Register",
				Name:  "register",
			},
		},
		help:    help.NewRegister(),
		actions: acts,
		log:     log,
	}

	return &m
}

func (m *RegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RegisterModel) lastIndex() int {
	return len(m.inputs) + len(m.buttons) - 1
}

func (m *RegisterModel) getInputValues() models.ClientAuthValues {
	return models.ClientAuthValues{
		Login:          m.inputs[0].Value(),
		Password:       m.inputs[1].Value(),
		CryptoPassword: m.inputs[2].Value(),
	}
}

func (m *RegisterModel) clearErrors() {
	m.errors = registerViewErrors{}
}

func (m *RegisterModel) changeCurrentIndex(keyMsg tea.KeyMsg) {
	if keyMsg.Type == tea.KeyUp || keyMsg.Type == tea.KeyShiftTab {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > m.lastIndex() {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = m.lastIndex()
	}
}

func (m *RegisterModel) ChangeFocus(keyMsg tea.KeyMsg) []tea.Cmd {
	m.changeCurrentIndex(keyMsg)

	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	for i := range m.buttons {
		btnIndex := len(m.inputs) + i

		if btnIndex == m.focusIndex {
			m.buttons[i].Focus()
		} else {
			m.buttons[i].Blur()
		}
	}

	return cmds
}

func (m *RegisterModel) isNavigationKey(key tea.KeyType) bool {
	return key == tea.KeyTab || key == tea.KeyShiftTab ||
		key == tea.KeyUp || key == tea.KeyDown || key == tea.KeyEnter
}

func (m *RegisterModel) handleKeyPress(keyMsg tea.KeyMsg) tea.Cmd {
	key := keyMsg.Type

	if !m.isNavigationKey(key) {
		return nil
	}

	if key == tea.KeyEnter && m.focusIndex >= len(m.inputs) {
		btn := m.buttons[m.focusIndex-len(m.inputs)]
		if btn.Name == "register" {
			return func() tea.Msg {
				m.clearErrors()

				res, err := m.actions.Register(context.Background(), m.getInputValues())
				if err != nil {
					return msgs.RegisterErrorMsg{Err: err}
				}

				return msgs.RegisterSuccessMsg{Res: res}
			}
		}

		return nil
	}

	cmds := m.ChangeFocus(keyMsg)

	return tea.Batch(cmds...)
}

func (m *RegisterModel) addErrors(msg msgs.RegisterErrorMsg) tea.Cmd {
	var resErr *req.ResError

	if errors.As(msg.Err, &resErr) && resErr != nil && resErr.Errors != nil {
		for key, val := range resErr.Errors {
			switch key {
			case "login":
				m.errors.Login = val
			case "password":
				m.errors.Password = val
			default:
				if val != "" {
					m.errors.Message += fmt.Sprintf("%s:%s ", key, val)
				}
			}
		}
	} else {
		m.errors.Message = msg.Err.Error()
	}

	return nil
}

func (m *RegisterModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds := []tea.Cmd{
			m.updateInputs(msg),
			m.handleKeyPress(msg),
		}

		return tea.Batch(cmds...)
	case msgs.RegisterErrorMsg:
		return m.addErrors(msg)
	default:
		return nil
	}
}

func (m *RegisterModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *RegisterModel) View() string {
	var b strings.Builder

	for i := 0; i < len(m.inputs); i++ {
		b.WriteString(m.inputs[i].View())

		switch m.inputs[i].Name {
		case "login":
			b.WriteString(styles.ErrorStyle.Render(m.errors.Login))
		case "password":
			b.WriteString(styles.ErrorStyle.Render(m.errors.Password))
		case "crypto-password":
			b.WriteString(styles.ErrorStyle.Render(m.errors.CryptoPassword))
		}

		b.WriteRune('\n')
	}

	fmt.Fprintf(&b, "%s\n\n", styles.ErrorStyle.Render(m.errors.Message))

	for i := 0; i < len(m.buttons); i++ {
		b.WriteString(m.buttons[i].View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
