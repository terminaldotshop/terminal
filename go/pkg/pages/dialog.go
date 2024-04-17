package pages

import "github.com/charmbracelet/lipgloss"

func DisplayDialog(m Model, display string) string {
    dialogBoxStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#874BFD")).
        Padding(1, 1, 1, 1).
        BorderTop(true).
        BorderLeft(true).
        BorderRight(true).
        BorderBottom(true)

    subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

    return lipgloss.Place(m.width, m.GetMaxPageHeight(),
        lipgloss.Center, lipgloss.Center,
        dialogBoxStyle.Render(display),
        lipgloss.WithWhitespaceChars("咖啡"),
        lipgloss.WithWhitespaceForeground(subtle),
    )
}

