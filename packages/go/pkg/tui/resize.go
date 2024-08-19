package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) ResizeView() string {
	return lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.theme.TextAccent().Render("your"),
			m.LogoView(),
			m.theme.TextAccent().Render("is too small"),
		),
	)
}
