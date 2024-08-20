package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminaldotshop/terminal/go/pkg/tui/theme"
)

type shopState struct {
	selected int
}

func (m model) ShopSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(shopPage)
	m.state.footer.commands = []footerCommand{
		{key: "+/-", value: "qty"},
		{key: "c", value: "cart"},
		{key: "q", value: "quit"},
	}

	if len(m.products) > 1 {
		m.state.footer.commands = append(
			[]footerCommand{{key: "↑↓", value: "products"}},
			m.state.footer.commands...,
		)
	}

	m = m.UpdateSelectedTheme()
	return m, nil
}

func (m model) ShopUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			return m.UpdateSelected(false)
		case "shift+tab", "up", "k":
			return m.UpdateSelected(true)
		case "+", "=", "right", "l":
			productVariantID := m.products[m.state.shop.selected].Variants[0].ID
			return m.UpdateCart(productVariantID, 1)
		case "-", "left", "h":
			productVariantID := m.products[m.state.shop.selected].Variants[0].ID
			return m.UpdateCart(productVariantID, -1)
		case "enter":
			return m.CartSwitch()
		}
	}

	return m, nil
}

func (m model) ShopView() string {
	base := m.theme.Base().Render
	accent := m.theme.TextAccent().Render
	bold := m.theme.TextHighlight().Bold(true).Render

	product := m.products[m.state.shop.selected]
	variantID := product.Variants[0].ID
	cartItem, _ := m.GetCartItem(variantID)

	minus := base("- ")
	plus := base(" +")
	count := accent(fmt.Sprintf(" %d ", cartItem.Quantity))
	quantity := minus + count + plus

	menuWidth := 0
	products := strings.Builder{}
	for _, product := range m.products {
		w := lipgloss.Width(product.Name)
		if w > menuWidth {
			menuWidth = w
		}
	}

	var menuItem lipgloss.Style
	var highlightedMenuItem lipgloss.Style

	if m.size < large {
		menuWidth = m.widthContent

		menuItem = m.theme.Base().
			Width(menuWidth).
			Align(lipgloss.Center)
		highlightedMenuItem = m.theme.Base().
			Width(menuWidth).
			Align(lipgloss.Center).
			Background(m.theme.Highlight()).
			Foreground(m.theme.Accent())
	} else {
		menuItem = m.theme.Base().
			Width(menuWidth+2).
			Padding(0, 1)
		highlightedMenuItem = m.theme.Base().
			Width(menuWidth+2).
			Padding(0, 1).
			Background(m.theme.Highlight()).
			Foreground(m.theme.Accent())
	}

	// TODO: do we need a category header?
	// products.WriteString(menuItem.Copy().Background(m.theme.Body()).Foreground(m.theme.Accent()).Render("coffee beans"))
	// products.WriteString("\n\n")

	for i := range m.products {
		name := m.products[i].Name

		var content string
		if i == m.state.shop.selected {
			content = highlightedMenuItem.Render(name)
		} else {
			content = menuItem.Render(name)
		}

		products.WriteString(content + "\n")
	}

	productList := m.theme.Base().Render(products.String())
	productListWidth := lipgloss.Width(productList)

	detailPaddingLeft := 2
	detailWidth := m.widthContent - productListWidth - detailPaddingLeft
	detailStyle := m.theme.Base().
		PaddingLeft(detailPaddingLeft).
		Width(detailWidth)
	name := accent(product.Name)
	// TODO: ratings? real ones?
	// nameWidth := lipgloss.Width(name)
	// rating := accent("★★★★★")
	// ratingWidth := lipgloss.Width(rating)
	// ratingSpace := m.theme.Base().Width(detailWidth - ratingWidth - nameWidth - 2).Render()

	detail := lipgloss.JoinVertical(
		lipgloss.Left,
		name, //+ratingSpace+rating,
		base(strings.ToLower(product.Variants[0].Name)),
		"",
		bold(fmt.Sprintf("$%.2v", product.Variants[0].Price/100)),
		"",
		product.Description,
		"\n",
		quantity,
	)

	var content string
	if len(m.products) == 1 {
		content = m.theme.Base().Width(m.widthContent).Render(detail)
	} else if m.size < large {
		detailStyle := m.theme.Base().
			Width(m.widthContent)

		content = m.theme.Base().
			Width(m.widthContent).
			Render(lipgloss.JoinVertical(
				lipgloss.Top,
				productList,
				detailStyle.Render(detail),
			))
	} else {
		content = m.theme.Base().
			Width(m.widthContent).
			Render(lipgloss.JoinHorizontal(
				lipgloss.Top,
				productList,
				detailStyle.Render(detail),
			))
	}

	return content
}

func (m model) UpdateSelectedTheme() model {
	var highlight string
	product := m.products[m.state.shop.selected]
	if strings.ToLower(product.Name) == "segfault" {
		highlight = "#169FC1"
	} else if strings.ToLower(product.Name) == "dark mode" {
		highlight = "#118B39"
	} else if strings.ToLower(product.Name) == "[object object]" {
		highlight = "#F5BB1D"
	} else if strings.ToLower(product.Name) == "404" {
		highlight = "#D53C81"
	} else if strings.ToLower(product.Name) == "artisan" {
		highlight = "#F9322C"
	}

	if highlight != "" {
		m.theme = theme.BasicTheme(m.renderer, &highlight)
	}

	return m
}

func (m model) UpdateSelected(previous bool) (model, tea.Cmd) {
	var next int
	if previous {
		next = m.state.shop.selected - 1
	} else {
		next = m.state.shop.selected + 1
	}

	if next < 0 {
		next = 0
	}
	max := len(m.products) - 1
	if next > max {
		next = max
	}

	m.state.shop.selected = next
	m = m.UpdateSelectedTheme()
	return m, nil
}
