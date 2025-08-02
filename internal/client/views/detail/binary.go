package detail

import (
	"strconv"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

func NewBinary(name string, val models.VaultItemDataBinary, acts *actions.Actions, id int) tea.Model {
	return components.NewDetail(
		[]components.Field{
			{
				Label: "Name",
				Value: name,
			},
			{
				Label: "Path",
				Value: val.Path,
			},
			{
				Label: "Size",
				Value: strconv.Itoa(int(val.Size)),
			},
		},
		acts,
		components.BinaryOpts{
			ID:   id,
			Path: val.Path,
		})
}
