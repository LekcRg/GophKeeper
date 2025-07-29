package views

import (
	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/client/router"
	"github.com/LekcRg/GophKeeper/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type Views struct {
	actions *actions.Actions
	log     *zap.Logger
	router  router.ViewRouter
}

func New(logger *zap.Logger, cfg *config.ClientConfig) *Views {
	if cfg == nil {
		cfg = &config.ClientConfig{}
	}

	request := req.New(cfg)
	acts := actions.New(request, logger, cfg)

	currentView := router.SelectAuthView
	if cfg.Key != "" {
		currentView = router.CryptoPassView
	}

	v := router.Views{
		router.SelectAuthView:  NewSelectAuth(cfg.Address),
		router.RegisterView:    NewRegister(acts, logger),
		router.TokenAuthView:   NewKeyAuth(acts),
		router.UpdateTokenView: NewUpdateKey(acts),
		router.CryptoPassView:  NewCryptoPass(acts),
	}

	m := &Views{
		router:  *router.NewViewRouter(currentView, v),
		actions: acts,
		log:     logger,
	}

	return m
}

func (m *Views) Init() tea.Cmd {
	return m.router.Init()
}

func (m *Views) handleUpdateCredentials(successMsg msgs.CredentialsBytesMsg) tea.Cmd {
	err := m.actions.UpdateConfigCredentials(successMsg)
	if err != nil {
		m.log.Error("Update credentials config err", zap.Error(err))

		return func() tea.Msg {
			return msgs.ErrorMsg(err)
		}
	}

	return m.router.SwitchTo(router.TokenAuthView)
}

func (m *Views) handleSelectAuth(msg msgs.SelectAuthMsg) tea.Cmd {
	err := m.actions.UpdateConfigAddress(msg.Address)
	if err != nil {
		m.log.Warn("Failed to persist address config, continuing anyway", zap.Error(err))

		return func() tea.Msg {
			return msgs.ErrorMsg(err)
		}
	}

	return m.router.SwitchTo(router.CurrentView(msg.Selected))
}

func (m *Views) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case tea.KeyMsg:
		if typeMsg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if typeMsg.Type == tea.KeyEsc && m.router.IsAuthenticationView() {
			return m, m.router.SwitchTo(router.SelectAuthView)
		}
	case msgs.SelectAuthMsg:
		return m, m.handleSelectAuth(typeMsg)
	case msgs.CredentialsBytesMsg:
		return m, m.handleUpdateCredentials(typeMsg)
	case msgs.UpdateKeySuccessMsg:
		return m, m.router.SwitchTo(router.SelectAuthView)
	case msgs.CryptoPassValid:
		return m, m.router.SwitchTo(router.ListView)
	}

	currentView := m.router.Current()
	if currentView != nil {
		newCurrentM, cmd := currentView.Update(msg)
		m.router.SetCurrentModel(newCurrentM)

		return m, cmd
	}

	return m, nil
}

func (m *Views) View() string {
	currentView := m.router.Current()
	if currentView == nil {
		m.log.Error("Not found current view")
		return "Error, ctrl+c to quit"
	}

	return currentView.View()
}
