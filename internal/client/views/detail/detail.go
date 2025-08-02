package detail

import (
	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/state"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type DetailModel struct {
	state   *state.State
	view    tea.Model
	log     *zap.Logger
	actions *actions.Actions
}

func NewDetail(st *state.State, log *zap.Logger, acts *actions.Actions) tea.Model {
	return &DetailModel{
		state:   st,
		log:     log,
		actions: acts,
	}
}

func (m *DetailModel) Init() tea.Cmd {
	item, err := m.state.GetActiveItem()
	if err != nil {
		m.log.Error("detail get item error", zap.Error(err))
		m.view = nil

		return nil
	}

	switch typedItem := item.DecryptedData.(type) {
	case models.VaultItemDataPassword:
		m.view = NewPassword(item.Name, typedItem, m.actions)
	case models.VaultNote:
		m.view = NewNote(item.Name, typedItem, m.actions)
	case models.VaultItemDataCard:
		m.view = NewCard(item.Name, typedItem, m.actions)
	case models.VaultItemDataBinary:
		m.view = NewBinary(item.Name, typedItem, m.actions, item.ID)
	}

	return func() tea.Msg { return "" }
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)

	return m, cmd
}

func (m *DetailModel) View() string {
	if m.view == nil {
		return "\nError retrieving the active vaul item, press ESC to back"
	}

	return m.view.View()
}
