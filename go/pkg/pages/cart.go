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
    productsContainer := m.renderer.NewStyle().
        Width(centeredWidth).
        Margin(0, 10, 0, 10)

    totalProductCount := m.theme.
        ActiveTitleForeground().
        Width(centeredWidth).
        Bold(true).
        Border(lipgloss.RoundedBorder(), false, false, true, false).
        MarginBottom(1).
        Align(lipgloss.Right).
        Render(fmt.Sprintf("Items: %d", m.cart.totalItems))

    productOrders := make([]string, 0)
    productOrders = append(productOrders, totalProductCount)

    for _, productInfo := range m.cart.products {

        count := fmt.Sprintf("Count: %d", productInfo.count)
        price := fmt.Sprintf("$%.2f", productInfo.product.Price * float64(productInfo.count))

        productCount := m.theme.
            NormalForeground().
            Bold(true).
            Render(count)

        title := m.theme.TitleForeground().
            Width(centeredWidth - len(count)).
            Render(productInfo.product.Name)

        description := m.theme.DescForeground().
            Width(centeredWidth).
            Render(productInfo.product.Description)

        priceTxt := m.theme.PriceForeground().
            Width(centeredWidth).
            AlignHorizontal(lipgloss.Right).
            MarginBottom(1).
            Render(price)

        productOrders = append(productOrders,
            lipgloss.JoinVertical(
                0,
                lipgloss.JoinHorizontal(0, title, productCount),
                lipgloss.JoinHorizontal(0, description),
                priceTxt,
            ),
        )
    }

    return productsContainer.Render(productOrders...)
}
*/
