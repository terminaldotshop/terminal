package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func cityStateZip(theme Theme, city, state, zip string) string {
    return fmt.Sprintf("%s, %s %s",
        theme.ActiveDescForeground().Render(city),
        theme.ActiveDescForeground().Render(state),
        theme.ActiveDescForeground().Render(zip))
}

func value(theme Theme, value string) string {
    return theme.ActiveDescForeground().Render(value)
}

func RenderShipping(m Model, shipping ShippingState, title string) string {
    theme := m.theme
    lines := []string{
        theme.ActiveTitleForeground().Render(title),
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

func expiration(theme Theme, month, year string) string {
    return theme.ActiveDescForeground().Render(fmt.Sprintf("%s/%s", month, year))
}

func cc(theme Theme, cc string) string {
    return theme.ActiveDescForeground().Render(cc[len(cc) - 4:])
}

func keyValue(theme Theme, key, value string) string {
    return fmt.Sprintf("%s: %s", theme.ActiveDescForeground().Render(key), theme.DescForeground().Render(value))
}


func RenderCreditCard(m Model, credit CreditCardState) string {
    theme := m.theme

    return lipgloss.JoinVertical(
        0,
        theme.ActiveTitleForeground().Render("Credit Card Info"),
        value(theme, credit.Name),
        cc(theme, credit.CC),
        lipgloss.JoinHorizontal(
            0,
            keyValue(theme, "CVC", "***"),
            expiration(theme, credit.ExpMonth, credit.ExpYear)))
}

func RenderEmail(m Model) string {
    theme := m.theme

    return lipgloss.JoinVertical(
        0,
        theme.ActiveTitleForeground().Render("Email"),
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
