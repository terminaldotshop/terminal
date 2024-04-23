package pages

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminalhq/terminal/go/pkg/api"
)

type ProductPage struct {
	Product *api.Product
}

func NewProductPage(product *api.Product) *ProductPage {
	return &ProductPage{
		Product: product,
	}
}

func (p *ProductPage) Exit(m Model) Model { return m }
func (s *ProductPage) Enter(m Model)      {}
func (w *ProductPage) Update(m Model, raw tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg := raw.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c", "C", "enter":
			if m.order.count == 0 {
				return true, m, NewDialog("You must select how many bags of coffee you wish to purchase")
			}
			return true, m, NewNavigateEmail
		case "j", "<", "left", "-":
			m.order.count = int(math.Max(0, float64(m.order.count-1)))
			return true, m, nil
		case "k", ">", "right", "+":
			m.order.count = int(math.Min(5, float64(m.order.count+1)))
			return true, m, nil
		}
	}
	return false, m, nil
}

func (w *ProductPage) Title() string { return "Order" }

func (w *ProductPage) Render(m *Model) string {
	titleStyle := m.theme.ActiveTitle().
		AlignHorizontal(lipgloss.Center).
		Width(m.width).
		MarginBottom(2)

	artWidth := m.width / 2
	descWidth := m.width - artWidth
	pageHeight := m.height - 5

	art := m.renderer.NewStyle().
		Width(artWidth)

	desc := m.theme.ActiveDescription().
		Width(descWidth).
		Render(m.order.product.Description)

	checkoutSpacing := m.renderer.NewStyle().
		Width(descWidth).
		MarginTop(pageHeight - lipgloss.Height(desc) - 10).
		Render("")

	countLeft := m.theme.HRB().
		Render(" - ")

	countRight := m.theme.HRB().
		Render(" + ")

	countStr := fmt.Sprintf("%d", m.order.count)
	count := m.renderer.NewStyle().
		Margin(1, 1, 0, 1).
		Render(countStr)

	countSpacing := m.renderer.NewStyle().
		Width(descWidth - (lipgloss.Width(countRight) + lipgloss.Width(countLeft) + lipgloss.Width(count))).
		Render("")

	row := lipgloss.JoinHorizontal(
		lipgloss.Left,
		art.Render(m.order.product.Art),
		lipgloss.JoinVertical(
			lipgloss.Top,
			desc,
			checkoutSpacing,
			lipgloss.JoinHorizontal(0, countLeft, count, countRight, countSpacing)))

	return lipgloss.JoinVertical(0,
		titleStyle.Render(fmt.Sprintf("===== %s =====", m.order.product.Name)),
		row,
	)
}
