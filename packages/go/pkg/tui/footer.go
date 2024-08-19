package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type footerState struct {
	commands []footerCommand
}

type footerCommand struct {
	key   string
	value string
}

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render

	table := m.theme.Base().
		Width(m.widthContainer).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.theme.Border()).
		PaddingBottom(1).
		Align(lipgloss.Center)

	if m.size == small && m.hasMenu {
		return table.Render(bold("m") + base(" menu"))
	}

	commands := []string{}
	for _, cmd := range m.state.footer.commands {
		commands = append(commands, bold(" "+cmd.key+" ")+base(cmd.value+"  "))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"free shipping on US orders over $40",
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				commands...,
			),
		))
}
