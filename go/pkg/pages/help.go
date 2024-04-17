package pages

import "github.com/charmbracelet/lipgloss"

func getHelpBasedOnPage(page int) string {
    switch page {
    case 0:
        return "You did too much enhance, please unenhance"
    }

    return ""
}

func helpMenu(m Model) string {
    helpContainer := lipgloss.NewStyle().
        Border(lipgloss.DoubleBorder(), true, false, false, false).
        Width(m.width).
        Height(2)

    return helpContainer.Render(getHelpBasedOnPage(m.currentPage))
}

