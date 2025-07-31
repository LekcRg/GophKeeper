package nav

import (
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Navigation struct {
	Up         key.Binding
	Down       key.Binding
	Inputs     []components.TextInput
	Buttons    []components.Button
	focusIndex int
}

func (n *Navigation) lastIndex() int {
	return len(n.Inputs) + len(n.Buttons) - 1
}

func (n *Navigation) MoveToNext() []tea.Cmd {
	n.focusIndex++
	if n.focusIndex > n.lastIndex() {
		n.focusIndex = 0
	}

	return n.updateFocus()
}

func (n *Navigation) MoveToPrev() []tea.Cmd {
	n.focusIndex--
	if n.focusIndex < 0 {
		n.focusIndex = n.lastIndex()
	}

	return n.updateFocus()
}

func (n *Navigation) updateFocus() []tea.Cmd {
	cmds := make([]tea.Cmd, len(n.Inputs))

	for i := range n.Inputs {
		if i == n.focusIndex {
			cmds[i] = n.Inputs[i].Focus()
		} else {
			n.Inputs[i].Blur()
		}
	}

	for i := range n.Buttons {
		btnIndex := len(n.Inputs) + i
		if btnIndex == n.focusIndex {
			n.Buttons[i].Focus()
		} else {
			n.Buttons[i].Blur()
		}
	}

	return cmds
}

func (n *Navigation) HandleKeyPress(k tea.KeyMsg) []tea.Cmd {
	switch {
	case key.Matches(k, n.Up):
		return n.MoveToPrev()
	case key.Matches(k, n.Down):
		return n.MoveToNext()
	}

	return nil
}

func (n *Navigation) IsOnButton() bool {
	return n.focusIndex >= len(n.Inputs)
}

func (n *Navigation) GetCurrentButton() *components.Button {
	if !n.IsOnButton() {
		return nil
	}

	return &n.Buttons[n.focusIndex-len(n.Inputs)]
}
