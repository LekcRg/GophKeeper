package detail

import (
	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

func NewNote(name string, val models.VaultNote, acts *actions.Actions) tea.Model {
	return components.NewDetail([]components.Field{
		{
			Label: "Name",
			Value: name,
		},
		{
			Label: "Content",
			Value: val.Text,
		},
	}, acts, components.BinaryOpts{})
}
