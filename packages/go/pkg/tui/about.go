package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) AboutSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(aboutPage)
	m.state.footer.commands = []footerCommand{
		{key: "c", value: "cart"},
	}
	return m, nil
}

func (m model) AboutUpdate(msg tea.Msg) (model, tea.Cmd) {
	return m, nil
}

func (m model) AboutView() string {
	base := m.theme.Base().Width(m.widthContent).Render
	accent := m.theme.TextAccent().Render

	return lipgloss.JoinVertical(
		lipgloss.Left,
		base("1. # Amazingly awesome products for developers brought to you by a group of talented, good looking, and humble heroes..."),
		"",
		base("2. # @thdxr"),
		"",
		base("3. # @adamdotdev"),
		"",
		base("4. # @theprimeagen"),
		"",
		base("5. # @teej_dv"),
		"",
		base("6. # @iamdavidhill"),
		"",
		m.theme.Base().Render("7. ")+accent("Terminal Products, Inc.")+m.CursorView(),
	)
}
