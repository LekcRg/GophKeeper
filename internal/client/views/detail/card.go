package detail

import (
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

func NewCard(name string, val models.VaultItemDataCard) tea.Model {
	return components.NewDetail([]components.Field{
		{
			Label: "Name",
			Value: name,
		},
		{
			Label: "Card number",
			Value: val.Number,
		},
		{
			Label: "Expire",
			Value: val.Exp,
		},
		{
			Label: "CVV",
			Value: val.CVV,
		},
	})
}
