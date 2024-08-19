package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	terminal "github.com/terminaldotshop/terminal-sdk-go"
)

func (m model) FinalSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(finalPage)
	m.state.footer.commands = []footerCommand{
		{key: "enter", value: "done"},
	}
	m.cart.Items = []terminal.CartItem{}
	m.cart.Subtotal = 0
	return m, nil
}

func (m model) FinalUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) FinalView() string {
	return m.theme.Base().Width(m.widthContent).Render(lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.TextAccent().Render("Thank you for ordering with Terminal Products, Inc.")),
		"\n\nAt this very moment as you sit, stunned and in awe of the CLI experience that just befell you, a personalised order confirmation email is on its way to your inbox.\n\nSimultaneously, news of your order is being celebrated wildly by the team. Perhaps too wildly by some. Once the excitement of your order has subsided to manageable levels your order will be sealed, shipped, and tracked courtesy of our very own Chief of SST.\n\nYours sincerely,\n\nDax, Adam, Prime, Teej, David\n\nTerminal Products, Inc.",
		fmt.Sprintf("\n\nps. %s", m.theme.TextHighlight().Render("https://www.terminal.shop/xxx")),
	)
}
