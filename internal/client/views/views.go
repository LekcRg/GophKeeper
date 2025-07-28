package views

import (
	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/client/router"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/errs"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type Views struct {
	router  router.ViewRouter
	actions *actions.Actions
	log     *zap.Logger
	// view    router.CurrentView
}

func New(logger *zap.Logger, cfg *config.ClientConfig) *Views {
	request := req.New(cfg)
	acts := actions.New(request, logger, cfg)

	currentView := router.SelectAuthView
	if cfg != nil && cfg.Key != "" {
		currentView = router.ListView
	}

	addr := ""
	if cfg != nil && cfg.Address != "" {
		addr = cfg.Address
	}

	v := router.Views{
		router.SelectAuthView: NewSelectAuth(addr),
		router.RegisterView:   NewRegister(acts, logger),
	}

	m := &Views{
		router: *router.NewViewRouter(currentView, v),
		// view:    currentView,
		actions: acts,
		log:     logger,
	}

	return m
}

func (m *Views) Init() tea.Cmd {
	return nil
}

func (m *Views) successRegister(msg tea.Msg) tea.Cmd {
	successMsg, ok := msg.(msgs.RegisterSuccessMsg)
	if ok {
		err := m.actions.UpdateConfigCredentials(successMsg.Res)
		if err != nil {
			m.log.Error("Update credentials config err", zap.Error(err))

			return func() tea.Msg {
				return msgs.ErrorMsg(err)
			}
		}

		m.router.SwitchTo(router.TokenAuthView)

		return nil
	}

	m.log.Error("successRegister type is not msgs.RegisterSuccessMsg")

	return func() tea.Msg {
		return msgs.ErrorMsg(errs.ErrInvalidType)
	}
}

func (m *Views) selectAuthView(msg msgs.SelectAuthMsg) {
	err := m.actions.UpdateConfigAddress(msg.Address)
	if err != nil {
		m.log.Error("Update address config err", zap.Error(err))
	}

	switch msg.Selected {
	case "register":
		m.router.SwitchTo(router.RegisterView)
	case "token":
		m.router.SwitchTo(router.TokenAuthView)
	case "update-token":
		m.router.SwitchTo(router.UpdateTokenView)
	}
}

func (m *Views) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case tea.KeyMsg:
		if typeMsg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if typeMsg.Type == tea.KeyEsc && m.router.IsAuthenticationView() {
			m.router.SwitchTo(router.SelectAuthView)
			return m, nil
		}
	case msgs.SelectAuthMsg:
		m.selectAuthView(typeMsg)
		return m, nil
	case msgs.RegisterSuccessMsg:
		m.log.Info("!!!!!!!!!")
		return m, m.successRegister(typeMsg)
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
	if currentView != nil {
		return currentView.View()
	}

	return "Error, ctrl+c to quit"
}
