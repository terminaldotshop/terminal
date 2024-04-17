package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type WidgetPage struct {
	widget *shop.Widget
}

func (w *WidgetPage) Title() string { return "Widget" }

func (w *WidgetPage) Render(m *Model) string {
	titleStyle := m.renderer.NewStyle().
		Bold(true).Foreground(lipgloss.Color("#b294bb")).Underline(true).AlignHorizontal(lipgloss.Center).
		Width(m.width).
		MarginBottom(2)

	artHeight := len(strings.Split(w.widget.Art, "\n"))
	descriptionStyle := m.renderer.NewStyle().
		Margin(0, 2).Padding(0, 2).Height(artHeight)

	rightSide := descriptionStyle.Render(w.widget.Description)
	row := lipgloss.JoinHorizontal(lipgloss.Left, w.widget.Art, rightSide)

	return fmt.Sprintf(`%s
%s`, titleStyle.Render(fmt.Sprintf("===== %s =====", w.widget.Name)), row)
}
