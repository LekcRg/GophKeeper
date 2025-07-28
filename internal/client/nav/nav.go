package nav

import (
	"github.com/LekcRg/GophKeeper/internal/client/components"
	tea "github.com/charmbracelet/bubbletea"
)

type Navigation struct {
	focusIndex int
	Inputs     []components.TextInput
	Buttons    []components.Button
}

func (n *Navigation) lastIndex() int {
	return len(n.Inputs) + len(n.Buttons) - 1
}

func (n *Navigation) moveToNext() {
	n.focusIndex++
	if n.focusIndex > n.lastIndex() {
		n.focusIndex = 0
	}
}

func (n *Navigation) moveToPrev() {
	n.focusIndex--
	if n.focusIndex < 0 {
		n.focusIndex = n.lastIndex()
	}
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

func (n *Navigation) HandleKeyPress(keyMsg tea.KeyMsg) []tea.Cmd {
	switch keyMsg.Type {
	case tea.KeyUp, tea.KeyShiftTab:
		n.moveToPrev()
		return n.updateFocus()
	case tea.KeyDown, tea.KeyTab:
		n.moveToNext()
		return n.updateFocus()
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
