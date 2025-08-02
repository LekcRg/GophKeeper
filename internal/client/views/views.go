package views

import (
	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/client/router"
	"github.com/LekcRg/GophKeeper/internal/client/state"
	"github.com/LekcRg/GophKeeper/internal/client/views/auth"
	"github.com/LekcRg/GophKeeper/internal/client/views/create"
	"github.com/LekcRg/GophKeeper/internal/client/views/detail"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type Views struct {
	actions *actions.Actions
	log     *zap.Logger
	state   *state.State
	router  router.ViewRouter
}

func New(logger *zap.Logger, cfg *config.ClientConfig) *Views {
	if cfg == nil {
		cfg = &config.ClientConfig{}
	}

	request := req.New(cfg)
	state := state.New(request, cfg)
	acts := actions.New(request, logger, cfg, state)

	currentView := router.SelectAuthView
	if cfg.Key != "" {
		currentView = router.CryptoPassView
	}

	v := router.Views{
		router.SelectAuthView:      auth.NewSelect(cfg.Address),
		router.RegisterView:        auth.NewRegister(acts, logger),
		router.TokenAuthView:       auth.NewKey(acts),
		router.UpdateTokenView:     auth.NewUpdateKey(acts),
		router.CryptoPassView:      auth.NewCryptoPass(acts),
		router.ListView:            NewList(logger, state),
		router.SelectVaultType:     create.NewSelectType(),
		router.CreateVaultPassword: create.NewPassword(acts, logger),
		router.CreateVaultNote:     create.NewNote(acts, logger),
		router.CreateVaultCard:     create.NewCard(acts, logger),
		router.CreateVaultBinary:   create.NewBinary(acts, logger),
		router.Detail:              detail.NewDetail(state, logger),
		router.FilePicker:          create.NewFilePicker(),
	}

	return &Views{
		router:  *router.NewViewRouter(currentView, v),
		actions: acts,
		log:     logger,
		state:   state,
	}
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

	return m.router.SwitchTo(router.CryptoPassView)
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

func (m *Views) handleUpdateVaultState(msg msgs.CreateVaultSuccess) tea.Cmd {
	return func() tea.Msg {
		m.state.AddVaultItem(msg.Item)

		return msgs.UpdateAndSwitchToTable{}
	}
}

func (m *Views) handleBack() tea.Msg {
	switch {
	case m.router.IsAuthenticationView():
		return m.router.SwitchTo(router.SelectAuthView)
	case m.router.IsListBack():
		return m.router.SwitchTo(router.ListView)
	case m.router.IsCreateView():
		return m.router.SwitchTo(router.SelectVaultType)
	case m.router.CurrentViewRoute() == router.FilePicker:
		return m.router.SwitchTo(router.CreateVaultBinary)
	}

	return nil
}

func (m *Views) handleFileSelected(msg msgs.FilepickerSelected) tea.Msg {
	m.router.SwitchTo(router.CreateVaultBinary)

	return func() tea.Msg {
		return msg
	}
}

func (m *Views) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(typeMsg, help.Quit) {
			return m, tea.Quit
		}

		if key.Matches(typeMsg, help.Back) {
			return m, m.handleBack
		}
	case msgs.Back:
		return m, m.handleBack
	case msgs.SelectAuthMsg:
		return m, m.handleSelectAuth(typeMsg)
	case msgs.CredentialsBytesMsg:
		return m, m.handleUpdateCredentials(typeMsg)
	case msgs.UpdateKeySuccessMsg:
		return m, m.router.SwitchTo(router.SelectAuthView)
	case msgs.CryptoPassValid:
		m.state.SaveCryptoPassword(typeMsg)
		return m, m.router.SwitchTo(router.ListView)
	case msgs.ToCreateVaultItem:
		return m, m.router.SwitchTo(router.SelectVaultType)
	case msgs.SelectTypeMsg:
		return m, m.router.SwitchTo(router.CurrentView(typeMsg.Selected))
	case msgs.CreateVaultSuccess:
		return m, m.handleUpdateVaultState(typeMsg)
	case msgs.UpdateAndSwitchToTable:
		return m, tea.Batch(
			m.router.SwitchTo(router.ListView),
			func() tea.Msg { return msgs.ListLoaded{} },
		)
	case msgs.SelectVaultItem:
		m.state.SetActiveID(int(typeMsg))
		return m, m.router.SwitchTo(router.Detail)
	case msgs.FilepickerSelected:
		m.handleFileSelected(typeMsg)
	case msgs.OpenFilePicker:
		return m, m.router.SwitchTo(router.FilePicker)
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
