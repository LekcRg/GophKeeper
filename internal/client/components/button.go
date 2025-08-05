package components

import (
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type Button struct {
	Label          string
	Name           string
	AdditionalText string
	Focused        bool
}

func (b *Button) Update(_ tea.Msg) tea.Cmd {
	return nil
}

func (b *Button) View() string {
	var res strings.Builder

	if b.Focused {
		res.WriteString(styles.FocusedStyle.Render(fmt.Sprintf("[ %s ]", b.Label)))
	} else {
		res.WriteString(styles.BlurredStyle.Render(fmt.Sprintf("[ %s ]", b.Label)))
	}

	if b.AdditionalText != "" {
		res.WriteString(" " + styles.BlurredStyle.Italic(true).Render(b.AdditionalText))
	}

	return res.String()
}

func (b *Button) Focus() {
	b.Focused = true
}

func (b *Button) Blur() {
	b.Focused = false
}
