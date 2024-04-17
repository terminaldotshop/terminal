package pages

import (
	"fmt"
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

    art := lipgloss.NewStyle().
        Width(artWidth)

    desc := lipgloss.NewStyle().
        Width(descWidth)

	row := lipgloss.JoinHorizontal(
        lipgloss.Left,
        art.Render(m.order.widget.Art),
        desc.Render(m.order.widget.Description))

	return fmt.Sprintf(`%s
%s`, titleStyle.Render(fmt.Sprintf("===== %s =====", m.order.widget.Name)), row)
}
