package views

import (
	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type CurrentView int

const (
	selectAuth CurrentView = iota
	register
	tokenAuth
	updateToken
	list
)

type Views struct {
	// tea.Model
	selectAuth   *SelectAuthModel
	register     *RegisterModel
	view         CurrentView
	securityData models.ClientRegisterResponse
}

func New(logger *zap.Logger) *Views {
	request := req.New()
	acts := actions.New(request, logger)
	m := &Views{
		register:   NewAuth(acts, logger),
		selectAuth: NewSelectAuth(),
		view:       selectAuth,
	}

	return m
}

func (m *Views) Init() tea.Cmd {
	return nil
}

func (m *Views) registerUpdate(msg tea.Msg) tea.Cmd {
	successMsg, ok := msg.(msgs.RegisterSuccessMsg)
	if ok {
		m.securityData = successMsg.Res
		m.view = tokenAuth
	}

	cmd := m.register.Update(msg)

	return cmd
}

func (m *Views) selectAuthView(msg msgs.SelectAuthMsg) {
	view := string(msg)
	switch view {
	case "register":
		m.view = register
	case "token":
		m.view = tokenAuth
	case "update-token":
		m.view = updateToken
	}
}

func (m *Views) isAuthenticationView() bool {
	return m.view == register || m.view == tokenAuth || m.view == updateToken
}

func (m *Views) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case tea.KeyMsg:
		if typeMsg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if typeMsg.Type == tea.KeyEsc && m.isAuthenticationView() {
			m.view = selectAuth
			return m, nil
		}
	case msgs.SelectAuthMsg:
		m.selectAuthView(typeMsg)
		return m, nil
	case msgs.RegisterSuccessMsg:
		return m, nil
	}

	switch m.view {
	case selectAuth:
		return m, m.selectAuth.Update(msg)
	case register:
		return m, m.registerUpdate(msg)
	}

	return m, nil
}

func (m *Views) View() string {
	switch m.view {
	case selectAuth:
		return m.selectAuth.View()
	case register:
		return m.register.View()
	}

	return "Error, ctrl+c to quit"
}
