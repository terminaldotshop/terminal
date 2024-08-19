package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	terminal "github.com/terminaldotshop/terminal-sdk-go"
)

type cartState struct {
	selected     int
	lastUpdateID int64
}

type CartUpdatedMsg struct {
	updateID int64
	updated  terminal.Cart
}

func (m model) IsCartEmpty() bool {
	return m.CartItemCount() == 0
}

func (m model) CartItemCount() int {
	return len(m.VisibleCartItems())
}

func (m model) VisibleCartItems() []terminal.CartItem {
	items := []terminal.CartItem{}
	for _, i := range m.cart.Items {
		if i.Quantity > 0 {
			items = append(items, i)
		}
	}
	return items
}

func (m model) GetCartItem(productVariantID string) (terminal.CartItem, int) {
	index := -1
	for i, item := range m.cart.Items {
		if item.ProductVariantID == productVariantID {
			index = i
			break
		}
	}

	var cartItem terminal.CartItem
	if index == -1 {
		cartItem = terminal.CartItem{ProductVariantID: productVariantID, Quantity: 0, Subtotal: 0}
	} else {
		cartItem = m.cart.Items[index]
	}

	return cartItem, index
}

func (m model) GetProduct(cartItem terminal.CartItem) (*terminal.Product, int) {
	index := -1
	for i, product := range m.products {
		if product.Variants[0].ID == cartItem.ProductVariantID {
			index = i
			break
		}
	}

	var product *terminal.Product
	if index == -1 {
		return nil, index
	} else {
		product = &m.products[index]
	}

	return product, index
}

func (m model) CalculateSubtotal() int64 {
	subtotal := int64(0)
	for _, item := range m.cart.Items {
		for _, product := range m.products {
			variant := product.Variants[0]
			if variant.ID == item.ProductVariantID {
				subtotal += item.Quantity * variant.Price
			}
		}
	}
	return subtotal
}

func (m model) UpdateCart(productVariantID string, offset int64) (model, tea.Cmd) {
	cartItem, index := m.GetCartItem(productVariantID)
	product, _ := m.GetProduct(cartItem)

	next := cartItem.Quantity + offset
	if next < 0 {
		return m, nil
	}
	if index == -1 {
		cartItem.Quantity = next
		cartItem.Subtotal = product.Variants[0].Price * next
		m.cart.Items = append(m.cart.Items, cartItem)
	} else {
		m.cart.Items[index].Quantity = next
		m.cart.Items[index].Subtotal = product.Variants[0].Price * next
	}

	updateID := time.Now().UTC().UnixMilli()
	m.cart.Subtotal = m.CalculateSubtotal()
	m.state.cart.lastUpdateID = updateID

	return m, func() tea.Msg {
		params := terminal.CartSetItemParams{
			ProductVariantID: terminal.String(cartItem.ProductVariantID),
			Quantity:         terminal.Int(next),
		}
		response, err := m.client.Cart.SetItem(m.context, params)
		if err != nil {
			// log.Error(err)
		}
		return CartUpdatedMsg{
			updateID: updateID,
			updated:  response.Result,
		}
	}
}

func (m model) UpdateSelectedCartItem(previous bool) (model, tea.Cmd) {
	if m.IsCartEmpty() {
		return m, nil
	}

	var next int
	if previous {
		next = m.state.cart.selected - 1
	} else {
		next = m.state.cart.selected + 1
	}

	if next < 0 {
		next = 0
	}

	max := m.CartItemCount() - 1
	if next > max {
		next = max
	}

	m.state.cart.selected = next
	return m, nil
}

func (m model) CartSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(cartPage)
	m.state.footer.commands = []footerCommand{
		{key: "esc", value: "back"},
		{key: "↑/↓", value: "items"},
		{key: "+/-", value: "qty"},
		{key: "c", value: "checkout"},
	}

	return m, nil
}

func (m model) CartUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "tab":
			return m.UpdateSelectedCartItem(false)
		case "k", "up", "shift+tab":
			return m.UpdateSelectedCartItem(true)
		case "+", "=", "right", "l":
			if m.IsCartEmpty() {
				return m, nil
			}
			productVariantID := m.VisibleCartItems()[m.state.cart.selected].ProductVariantID
			return m.UpdateCart(productVariantID, 1)
		case "-", "left", "h":
			if m.IsCartEmpty() {
				return m, nil
			}
			productVariantID := m.VisibleCartItems()[m.state.cart.selected].ProductVariantID
			return m.UpdateCart(productVariantID, -1)
		case "enter", "c":
			if m.IsCartEmpty() {
				return m, nil
			}
			return m.PaymentSwitch()
		case "esc":
			return m.ShopSwitch()
		}
	}

	return m, nil
}

func (m model) CartView() string {
	base := m.theme.Base().Align(lipgloss.Left).Render
	accent := m.theme.TextAccent().Render

	if m.IsCartEmpty() {
		return lipgloss.Place(
			m.widthContent,
			m.heightContent,
			lipgloss.Center,
			lipgloss.Center,
			base("Your cart is empty."),
		)
	}

	var lines []string
	for i, item := range m.VisibleCartItems() {
		product, _ := m.GetProduct(item)
		name := accent(product.Name)
		description := base(strings.ToLower(product.Variants[0].Name))
		quantity := base("  ") + accent(strconv.FormatInt(item.Quantity, 10)) + base("    ")
		if m.state.cart.selected == i {
			quantity = base("- ") + accent(strconv.FormatInt(item.Quantity, 10)) + base(" +  ")
		}
		subtotal := m.theme.Base().Width(5).Render(fmt.Sprintf("$%v", item.Subtotal/100))
		space := m.widthContent - lipgloss.Width(
			name,
		) - lipgloss.Width(
			quantity,
		) - lipgloss.Width(
			subtotal,
		) - 4

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				name,
				m.theme.Base().Width(space).Render(),
				quantity,
				subtotal,
			),
			description,
		)

		line := m.CreateBox(content, i == m.state.cart.selected)
		lines = append(lines, line)
	}

	return m.theme.Base().Render(lipgloss.JoinVertical(
		lipgloss.Left,
		lines...,
	))
}
