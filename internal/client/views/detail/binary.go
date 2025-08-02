package detail

import (
	"strconv"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

func NewBinary(name string, val models.VaultItemDataBinary) tea.Model {
	return components.NewDetail([]components.Field{
		{
			Label: "Name",
			Value: name,
		},
		{
			Label: "Name",
			Value: val.Name,
		},
		{
			Label: "Size",
			Value: strconv.Itoa(int(val.Size)),
		},
	})
}
