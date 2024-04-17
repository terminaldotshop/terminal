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

