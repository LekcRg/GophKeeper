package form

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/nav"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Form struct {
	Errors     *Errors
	validRules []*validation.KeyRules
	nav        nav.Navigation
}

type Opts struct {
	Inputs    []components.TextInput
	Buttons   []components.Button
	Textareas []components.Textarea
	Up        key.Binding
	Down      key.Binding
}

func NewForm(opts Opts) *Form {
	inputNames := make([]string, len(opts.Inputs))
	validRules := make([]*validation.KeyRules, len(opts.Inputs))

	for i := range opts.Inputs {
		input := &opts.Inputs[i]
		inputNames[i] = input.Name
		validRules[i] = validation.Key(input.Name, input.Valid...)
	}

	m := &Form{
		nav: nav.Navigation{
			Inputs:    opts.Inputs,
			Buttons:   opts.Buttons,
			Textareas: opts.Textareas,
			Up:        opts.Up,
			Down:      opts.Down,
		},
		Errors:     NewErrors(inputNames),
		validRules: validRules,
	}

	return m
}

func (m *Form) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (m *Form) GetValues() map[string]string {
	res := make(map[string]string, len(m.nav.Inputs))

	for i := range m.nav.Inputs {
		input := &m.nav.Inputs[i]
		res[input.Name] = input.Value()
	}

	return res
}

func (m *Form) HandleError(err error, apiKey string) {
	m.Errors.HandleError(err, apiKey)
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

func (m *Form) UpdateTextarea(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.nav.Textareas))
	for i := range m.nav.Textareas {
		cmds[i] = m.nav.Textareas[i].Update(msg)
	}

	return cmds
}

func (m *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	inputCmds := append(m.UpdateTextarea(msg), m.updateInputs(msg)...)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		navCmds := m.nav.HandleKeyPress(msg)

		if key.Matches(msg, help.Select) {
			if m.nav.IsOnInputs() {
				navCmds = append(navCmds, m.nav.MoveToNext()...)

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

	if len(m.nav.Inputs) > 0 {
		for i := range m.nav.Inputs {
			input := &m.nav.Inputs[i]
			b.WriteString(input.View())
			b.WriteString(styles.ErrorStyle.Render(
				m.Errors.GetFieldError(input.Name),
			))
			b.WriteRune('\n')
		}
	}

	if len(m.nav.Inputs) > 0 && len(m.nav.Textareas) > 0 {
		b.WriteRune('\n')
	}

	for i := range m.nav.Textareas {
		ta := &m.nav.Textareas[i]
		b.WriteString(ta.View())
		b.WriteString(styles.ErrorStyle.Render(
			m.Errors.GetFieldError(ta.Name),
		))
		b.WriteRune('\n')
	}

	if len(m.nav.Inputs) > 0 || len(m.nav.Textareas) > 0 {
		b.WriteRune('\n')
		b.WriteString(styles.ErrorStyle.Render(m.Errors.Message))
		b.WriteRune('\n')
	}

	for _, btn := range m.nav.Buttons {
		b.WriteString(btn.View())
		b.WriteRune('\n')
	}

	return b.String()
}
