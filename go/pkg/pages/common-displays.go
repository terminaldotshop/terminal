package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/terminalhq/terminal/go/pkg/api"
)

func cityStateZip(theme Theme, city, state, zip string) string {
	return theme.Description().Render(fmt.Sprintf("%s, %s %s", city, state, zip))
}

func value(theme Theme, value string) string {
	return theme.Description().Render(value)
}

func RenderShipping(theme Theme, shipping api.Address, title string) string {
	lines := []string{
		theme.ActiveTitle().Render(title),
		value(theme, shipping.Name),
		value(theme, shipping.AddrLine1),
	}

	if len(strings.TrimSpace(shipping.AddrLine2)) > 0 {
		lines = append(lines, value(theme, shipping.AddrLine2))
	}

	lines = append(lines, cityStateZip(theme, shipping.City, shipping.State, shipping.Zip))

	return lipgloss.JoinVertical(
		0,
		lines...)
}

func RenderSameShipping(theme Theme, title string) string {
	return lipgloss.JoinVertical(
		0,
		theme.ActiveTitle().Render(title),
		theme.Description().Render("(Same as Shipping Address)"),
	)
}

func expiration(theme Theme, month, year string) string {
	return theme.Description().Render(fmt.Sprintf("Expires: %s/%s", month, year))
}

func cc(theme Theme, cc string) string {
	return theme.Description().Render(
		fmt.Sprintf("•••• •••• •••• %s", cc[len(cc)-4:]),
	)
}

func keyValue(theme Theme, key, value string) string {
	return fmt.Sprintf("%s: %s", theme.Description().Render(key), theme.Description().Render(value))
}

func RenderCreditCard(theme Theme, credit api.CreditCard) string {
	return lipgloss.JoinVertical(
		0,
		theme.ActiveTitle().Render("Credit Card Info"),
		value(theme, credit.Name),
		cc(theme, credit.Number),
		expiration(theme, credit.ExpMonth, credit.ExpYear),
	)
}

func RenderOrder(m Model, order *api.OrderResponse) string {
	return m.theme.HRB().Render(
		fmt.Sprintf("PRESS ENTER TO CONFIRM\nTotal Cost: $%d.%d", order.Total/100, order.Total%100),
	)
}

func RenderEmail(m Model) string {
	theme := m.theme

	return lipgloss.JoinVertical(
		0,
		theme.ActiveTitle().Render("Email"),
		value(theme, m.email))
}

func RenderSplitView(m Model, view1, view2 string) string {
	container := lipgloss.NewStyle().
		Height(m.GetMaxPageHeight())

	leftWidth := m.GetMaxPageWidth() / 2
	rightWidth := m.GetMaxPageWidth() - leftWidth

	leftContainer := lipgloss.NewStyle().
		Width(leftWidth)

	rightContainer := lipgloss.NewStyle().
		Width(rightWidth)

	return container.Render(
		lipgloss.JoinHorizontal(0,
			leftContainer.Render(view1),
			rightContainer.Render(view2)),
	)
}
