package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func helpNav(theme Theme, k, m string) string {
	key := theme.HelpSpecial().Render(k)
	message := theme.HelpDescription().Render(fmt.Sprintf(" = %s ", m))
	return lipgloss.JoinHorizontal(0, key, message)
}

func joinHelpItems(theme Theme, sections ...string) string {
	sep := theme.NonText().Render("-")

	items := make([]string, 0)
	for _, section := range sections {
		items = append(items, sep)
		items = append(items, section)
	}

	items = append(items, sep)
	return lipgloss.JoinHorizontal(0, items...)
}

func getHelpBasedOnPage(m Model) string {
	theme := m.theme
	if m.Dialog != nil {
		return joinHelpItems(theme, helpNav(m.theme, "<esc>", "close dialog"))
	}

	infoPageNav := []string{
		helpNav(m.theme, "shift+tab", "go back to product"),
		helpNav(m.theme, "enter", "next field or page"),
		helpNav(m.theme, "ctrl+c", "quit"),
	}

	switch m.currentPage {
	case MIN_WIDTH_NOT_MET_PAGE:
		return "Your terminal dimensions are too small"
	case PRODUCT_PAGE:
		return joinHelpItems(theme,
			helpNav(m.theme, "j", "remove one from your order count"),
			helpNav(m.theme, "k", "add one to your order count"),
			helpNav(m.theme, "c", "begin checkout"),
			helpNav(m.theme, "ctrl+c", "quit"),
		)
	case EMAIL_PAGE:
		fallthrough
	case SHIPPING_PAGE:
		fallthrough
	case CC_PAGE:
		fallthrough
	case CC_ADDR_PAGE:
		return joinHelpItems(theme, infoPageNav...)

	case CONFIRM_PAGE:
		return joinHelpItems(theme,
			helpNav(m.theme, "enter", "to confirm and place order"),
			helpNav(m.theme, "shift+tab", "to go back"),
			helpNav(m.theme, "ctrl+c", "quit"),
		)
	}

	return ""
}

func helpMenu(m Model) string {
	helpContainer := m.renderer.NewStyle().
		Border(lipgloss.DoubleBorder(), true, false, false, false).
		Width(m.width).
		Height(2)

	return helpContainer.Render(getHelpBasedOnPage(m))
}
