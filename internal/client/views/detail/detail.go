package detail

import (
	"github.com/LekcRg/GophKeeper/internal/client/state"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type DetailModel struct {
	state *state.State
	view  tea.Model
	log   *zap.Logger
}

func NewDetail(st *state.State, log *zap.Logger) tea.Model {
	return &DetailModel{
		state: st,
		log:   log,
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
		m.view = NewPassword(item.Name, typedItem)
		m.log.Info("1")
	case models.VaultNote:
		m.view = NewNote(item.Name, typedItem)
		m.log.Info("2")
	case models.VaultItemDataCard:
		m.view = NewCard(item.Name, typedItem)
		m.log.Info("3")
	case models.VaultItemDataBinary:
		m.view = NewBinary(item.Name, typedItem)
	}

	return func() tea.Msg { return "" }
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *DetailModel) View() string {
	if m.view == nil {
		return "\nError retrieving the active vaul item, press ESC to back"
	}

	return m.view.View()
}
