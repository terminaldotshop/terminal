package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/terminalhq/terminal/go/pkg/api"
)

type ConfirmPage struct{}

func NewConfirmPage() *ConfirmPage {
	return &ConfirmPage{}
}

var currentOrder *api.OrderResponse

func (c *ConfirmPage) Exit(m Model) Model { return m }
func (c *ConfirmPage) Enter(m Model) {
	// TODO: Need to update this...
	order, err := api.CreateOrder(m.userToken, api.OrderParams{
		Email:    m.email,
		Shipping: m.shippingAddress,
		Products: []api.ProductOrder{
			// TODO(launch): This is pretty ugly
			{
				ID:       m.order.product.ID,
				Quantity: m.order.count,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// TODO(launch): This will create a new order every time...
	currentOrder = order
}

func (s *ConfirmPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c", "C", "enter":
			success, err := api.PurchaseOrder(m.userToken, currentOrder.OrderID, m.stripeCardToken)
			if err != nil {
				return true, m, NewDialog(fmt.Sprintf("PURCHASE ERROR: %s", err))
			}

			if *success {
				return true, m, NewNavigateAnimation
			}

			return true, m, nil
		}
	}

	return false, m, nil
}

func (s *ConfirmPage) Title() string { return "Confirmation" }

func (s *ConfirmPage) Render(m *Model) string {
	container := lipgloss.NewStyle().
		Height(m.GetMaxPageHeight()).
		PaddingLeft(2)

	lines := []string{
		RenderEmail(*m),
		RenderShipping(m.theme, m.shippingAddress, "Shipping Address"),
		RenderCreditCard(m.theme, m.creditCard),
	}

	billingAddressTitle := "Billing Address"
	if m.differentBillingAddress {
		lines = append(lines, RenderShipping(m.theme, m.billingAddress, billingAddressTitle))
	} else {
		lines = append(lines, RenderSameShipping(m.theme, billingAddressTitle))
	}

	return container.Render(
		lipgloss.JoinVertical(
			0,
			lines...))
}
