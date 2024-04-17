package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func helpNav(theme Theme, k, m string) string {
    key := theme.HelpSpecialForeground().Render(k)
    message := theme.DescForeground().Render(fmt.Sprintf(" = %s", m))

    return lipgloss.JoinHorizontal(0, key, message)
}

func getHelpBasedOnPage(m Model) string {
    if m.Dialog != nil {
        return lipgloss.JoinHorizontal(0,
            helpNav(m.theme, "<esc>", "close dialog"),
        )
    }

    switch m.currentPage {
    case MIN_WIDTH_NOT_MET_PAGE:
        return "You did too much enhance, please unenhance"
    case PRODUCT_PAGE:
        return lipgloss.JoinHorizontal(0,
            helpNav(m.theme, "j", "remove one from your order count"),
            helpNav(m.theme, "k", "add one to your order count"),
            helpNav(m.theme, "c", "begin checkout"),
            helpNav(m.theme, "C-c", "quit"),
        )
    case EMAIL_PAGE:
        return lipgloss.JoinHorizontal(0,
            helpNav(m.theme, "S-tab", "go back to product"),
            helpNav(m.theme, "C-c", "quit"),
        )
    }

    return ""
}

func helpMenu(m Model) string {
    helpContainer := lipgloss.NewStyle().
        Border(lipgloss.DoubleBorder(), true, false, false, false).
        Width(m.width).
        Height(2)

    return helpContainer.Render(getHelpBasedOnPage(m))
}

