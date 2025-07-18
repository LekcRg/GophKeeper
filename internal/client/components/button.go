package components

import (
	"fmt"

	"github.com/LekcRg/GophKeeper/internal/client/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type Button struct {
	Label   string
	Focused bool
}

func (b *Button) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (b *Button) View() string {
	if b.Focused {
		return styles.FocusedStyle.Render(fmt.Sprintf("[ %s ]", b.Label))
	}

	return styles.BlurredStyle.Render(fmt.Sprintf("[ %s ]", b.Label))
}

func (b *Button) Focus() {
	b.Focused = true
}

func (b *Button) Blur() {
	b.Focused = false
}
