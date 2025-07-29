package form

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/nav"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Form struct {
	Errors     *Errors
	nav        nav.Navigation
	validRules []*validation.KeyRules
}

func NewForm(
	inputs []components.TextInput, buttons []components.Button,
) *Form {
	inputNames := make([]string, len(inputs))
	validRules := make([]*validation.KeyRules, len(inputs))

	for i := range inputs {
		input := &inputs[i]
		inputNames[i] = input.Name
		validRules[i] = validation.Key(input.Name, input.Valid...)
	}

	m := &Form{
		nav: nav.Navigation{
			Inputs:  inputs,
			Buttons: buttons,
		},
		Errors:     NewErrors(inputNames),
		validRules: validRules,
	}

	return m
}

func (m *Form) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Form) GetValues() map[string]string {
	res := make(map[string]string, len(m.nav.Inputs))

	for i := range m.nav.Inputs {
		input := &m.nav.Inputs[i]
		res[input.Name] = input.Value()
	}

	return res
}

func (m *Form) HandleError(err error, key string) {
	m.Errors.HandleError(err, key)
}

func (m *Form) Submit() tea.Msg {
	m.Errors.Clear()
	vals := m.GetValues()

	var err error
	if len(m.validRules) > 0 {
		err = valid.MapString(vals, m.validRules)
		if err != nil {
			return msgs.ErrorMsg(err)
		}
	}

	return msgs.FormSubmitMsg{
		Values:     vals,
		ButtonName: m.nav.GetCurrentButton().Name,
	}
}

func (m *Form) updateInputs(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.nav.Inputs))
	for i := range m.nav.Inputs {
		cmds[i] = m.nav.Inputs[i].Update(msg)
	}

	return cmds
}

func (m *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	inputCmds := m.updateInputs(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		navCmds := m.nav.HandleKeyPress(msg)

		if msg.Type == tea.KeyEnter {
			if !m.nav.IsOnButton() {
				m.nav.MoveToNext()
				return m, tea.Batch(append(inputCmds, navCmds...)...)
			}

			btn := m.nav.GetCurrentButton()
			if btn != nil {
				return m, m.Submit
			}
		}

		return m, tea.Batch(append(inputCmds, navCmds...)...)
	case msgs.ErrorMsg:
		m.Errors.HandleAPIError(msg)
	}

	return m, tea.Batch(inputCmds...)
}

func (m *Form) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	for i := range m.nav.Inputs {
		input := &m.nav.Inputs[i]
		b.WriteString(input.View())
		b.WriteString(styles.ErrorStyle.Render(
			m.Errors.GetFieldError(input.Name),
		))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(styles.ErrorStyle.Render(m.Errors.Message))
	b.WriteRune('\n')

	for _, btn := range m.nav.Buttons {
		b.WriteString(btn.View())
		b.WriteRune('\n')
	}

	return b.String()
}
