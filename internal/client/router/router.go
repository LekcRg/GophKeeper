package router

import (
	tea "github.com/charmbracelet/bubbletea"
)

type CurrentView string

type Views map[CurrentView]tea.Model

type ViewRouter struct {
	views       Views
	currentView CurrentView
}

const (
	SelectAuthView      CurrentView = "select-auth"
	RegisterView        CurrentView = "register"
	TokenAuthView       CurrentView = "token"
	UpdateTokenView     CurrentView = "update-token"
	CryptoPassView      CurrentView = "crypto-pass"
	ListView            CurrentView = "list"
	SelectVaultType     CurrentView = "select-vault-type"
	CreateVaultPassword CurrentView = "password"
	CreateVaultNote     CurrentView = "note"
	CreateVaultCard     CurrentView = "card"
	CreateVaultBinary   CurrentView = "binary"
)

func NewViewRouter(current CurrentView, v Views) *ViewRouter {
	return &ViewRouter{
		currentView: current,
		views:       v,
	}
}

func (r *ViewRouter) Init() tea.Cmd {
	current := r.Current()
	if current != nil {
		return current.Init()
	}

	return nil
}

func (r *ViewRouter) IsAuthenticationView() bool {
	return r.currentView == RegisterView ||
		r.currentView == TokenAuthView ||
		r.currentView == UpdateTokenView
}

func (r *ViewRouter) IsCreateView() bool {
	return r.currentView == CreateVaultPassword ||
		r.currentView == CreateVaultNote ||
		r.currentView == CreateVaultCard ||
		r.currentView == CreateVaultBinary
}

func (r *ViewRouter) SwitchTo(view CurrentView) tea.Cmd {
	r.currentView = view
	return r.Init()
}

func (r *ViewRouter) CurrentViewRoute() CurrentView {
	return r.currentView
}

func (r *ViewRouter) Current() tea.Model {
	return r.views[r.currentView]
}

func (r *ViewRouter) SetCurrentModel(m tea.Model) {
	r.views[r.currentView] = m
}
