package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type menuState struct {
	lastPage page
}

func (m model) MenuSwitch() (model, tea.Cmd) {
	m.state.menu.lastPage = m.page
	m = m.SwitchPage(menuPage)
	return m, nil
}

func (m model) MenuUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			switch m.state.menu.lastPage {
			case aboutPage:
				return m.AboutSwitch()
			case faqPage:
				return m.FaqSwitch()
			// case paymentPage:
			// 	return m.PaymentSwitch()
			// case cartPage:
			// 	return m.CartSwitch()
			default:
				return m.ShopSwitch()
			}
		}
	}

	return m, nil
}

func (m model) MenuView() string {
	base := m.theme.Base().Render
	bold := m.theme.TextAccent().Bold(true).Render

	menu :=
		table.New().
			Border(lipgloss.HiddenBorder()).
			Row(bold("s"), base("shop")).
			Row(bold("a"), base("about")).
			Row(bold("f"), base("faq")).
			Row(bold("c"), base("cart")).
			Row("").
			StyleFunc(func(row, col int) lipgloss.Style {
				return m.theme.Base().
					Padding(0, 1).
					AlignHorizontal(lipgloss.Left)
			})

	for _, cmd := range m.state.footer.commands {
		if cmd.key == "s" ||
			cmd.key == "a" ||
			cmd.key == "f" ||
			cmd.key == "c" {
			continue
		}

		menu.Row(bold(cmd.key), base(cmd.value))
	}

	modal := m.theme.Base().
		Padding(1).
		Border(lipgloss.NormalBorder(), true, false).
		BorderForeground(m.theme.Border()).
		Render

	return lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.LogoView(),
			modal(menu.Render()),
			m.theme.TextAccent().
				Width(m.widthContent).
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center).
				Render("press esc to close"),
		),
	)
}
