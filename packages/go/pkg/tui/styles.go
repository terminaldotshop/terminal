package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) _createBoxInner(content string, selected bool, position lipgloss.Position, padding int) string {
	totalWidth := m.widthContent - 2

	padded := lipgloss.PlaceHorizontal(totalWidth, position, content)
	base := m.theme.Base().Border(lipgloss.NormalBorder()).Width(totalWidth)

	var style lipgloss.Style
	if selected {
		style = base.BorderForeground(m.theme.Accent())
	} else {
		style = base.BorderForeground(m.theme.Border())
	}
	return style.PaddingLeft(padding).Render(padded)
}

func (m model) CreateBox(content string, selected bool) string {
	return m._createBoxInner(content, selected, lipgloss.Left, 1)
}

func (m model) CreateCenteredBox(content string, selected bool) string {
	return m._createBoxInner(content, selected, lipgloss.Center, 0)
}
