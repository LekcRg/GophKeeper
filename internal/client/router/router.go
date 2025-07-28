package router

import tea "github.com/charmbracelet/bubbletea"

type CurrentView string

type Views map[CurrentView]tea.Model

type ViewRouter struct {
	currentView CurrentView
	views       Views
}

const (
	SelectAuthView  CurrentView = "select-auth"
	RegisterView    CurrentView = "register"
	TokenAuthView   CurrentView = "token"
	UpdateTokenView CurrentView = "update-token"
	ListView        CurrentView = "list"
)

func NewViewRouter(current CurrentView, v Views) *ViewRouter {
	return &ViewRouter{
		currentView: current,
		views:       v,
	}
}

func (r *ViewRouter) IsAuthenticationView() bool {
	return r.currentView == RegisterView ||
		r.currentView == TokenAuthView ||
		r.currentView == UpdateTokenView
}

func (r *ViewRouter) SwitchTo(view CurrentView) {
	r.currentView = view
}

func (r *ViewRouter) Current() tea.Model {
	return r.views[r.currentView]
}

func (r *ViewRouter) SetCurrentModel(m tea.Model) {
	r.views[r.currentView] = m
}
