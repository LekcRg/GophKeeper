package detail

import (
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

func NewPassword(name string, val models.VaultItemDataPassword) tea.Model {
	return components.NewDetail([]components.Field{
		{
			Label: "Name",
			Value: name,
		},
		{
			Label: "Login",
			Value: val.Login,
		},
		{
			Label: "Password",
			Value: val.Password,
		},
		{
			Label: "Link",
			Value: val.URL,
		},
	})
}
