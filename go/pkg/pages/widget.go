package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type WidgetPage struct { }

func (w *WidgetPage) Title() string { return "Widget" }

func (w *WidgetPage) Render(m *Model) string {
	titleStyle := m.theme.ActiveTitleForeground().
        AlignHorizontal(lipgloss.Center).
		Width(m.width).
		MarginBottom(2)

    artWidth := m.width / 2
    descWidth := m.width - artWidth

    artContainer := lipgloss.NewStyle().
        Width(artWidth)

    descContainer := lipgloss.NewStyle().
        Width(descWidth)

	descriptionStyle := m.renderer.
        NewStyle().
		Margin(0, 2).
        Padding(0, 2)

	rightSide := descriptionStyle.Render(m.order.widget.Description)
	row := lipgloss.JoinHorizontal(lipgloss.Left, m.order.widget.Art, rightSide)

	return fmt.Sprintf(`%s
%s`, titleStyle.Render(fmt.Sprintf("===== %s =====", m.order.widget.Name)), row)
}
