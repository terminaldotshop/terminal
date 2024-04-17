package pages
/*

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type CartPage struct {
}

func NewCartPage(m Model) CartPage {
    return CartPage{ }
}

func (c *CartPage) Title() string { return "Cart" }

func (c *CartPage) Render(m *Model) string {
    centeredWidth := m.width - 20
    widgetsContainer := m.renderer.NewStyle().
        Width(centeredWidth).
        Margin(0, 10, 0, 10)

    totalWidgetCount := m.theme.
        ActiveTitleForeground().
        Width(centeredWidth).
        Bold(true).
        Border(lipgloss.RoundedBorder(), false, false, true, false).
        MarginBottom(1).
        Align(lipgloss.Right).
        Render(fmt.Sprintf("Items: %d", m.cart.totalItems))

    widgetOrders := make([]string, 0)
    widgetOrders = append(widgetOrders, totalWidgetCount)

    for _, widgetInfo := range m.cart.widgets {

        count := fmt.Sprintf("Count: %d", widgetInfo.count)
        price := fmt.Sprintf("$%.2f", widgetInfo.widget.Price * float64(widgetInfo.count))

        widgetCount := m.theme.
            NormalForeground().
            Bold(true).
            Render(count)

        title := m.theme.TitleForeground().
            Width(centeredWidth - len(count)).
            Render(widgetInfo.widget.Name)

        description := m.theme.DescForeground().
            Width(centeredWidth).
            Render(widgetInfo.widget.Description)

        priceTxt := m.theme.PriceForeground().
            Width(centeredWidth).
            AlignHorizontal(lipgloss.Right).
            MarginBottom(1).
            Render(price)

        widgetOrders = append(widgetOrders,
            lipgloss.JoinVertical(
                0,
                lipgloss.JoinHorizontal(0, title, widgetCount),
                lipgloss.JoinHorizontal(0, description),
                priceTxt,
            ),
        )
    }

    return widgetsContainer.Render(widgetOrders...)
}
*/
