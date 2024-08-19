package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.hasMenu {
			switch msg.String() {
			case "c":
				return m.CartSwitch()
			case "s":
				return m.ShopSwitch()
			case "a":
				return m.AboutSwitch()
			case "f":
				return m.FaqSwitch()
			case "m":
				return m.MenuSwitch()
			case "q":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) HeaderView() string {
	total := m.cart.Subtotal
	count := int64(0)
	if m.cart.Items != nil {
		for _, item := range m.cart.Items {
			count += item.Quantity
		}
	}

	bold := m.theme.TextAccent().Bold(true).Render
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render
	cursor := m.theme.Base().Background(m.theme.Highlight()).Render(" ")

	menu := bold("m") + base(" ☰")
	back := base("← ") + bold("esc") + base(" back")
	mark := bold("t") + cursor
	logo := bold("terminal")
	shop := accent("s") + base(" shop")
	about := accent("a") + base(" about")
	faq := accent("f") + base(" faq")
	cart :=
		accent("c") +
			base(" cart") +
			accent(fmt.Sprintf(" $%2v", total/100)) +
			base(fmt.Sprintf(" [%d]", count))

	switch m.page {
	case shopPage:
		shop = accent("s shop")
	case aboutPage:
		about = accent("a about")
	case faqPage:
		faq = accent("f faq")
	}

	var tabs []string

	switch m.size {
	case small:
		tabs = []string{
			mark,
			cart,
		}
	case medium:
		if m.hasMenu {
			tabs = []string{
				menu,
				logo,
				cart,
			}
		} else if m.checkout {
			tabs = []string{
				back,
				logo,
				cart,
			}
		} else {
			tabs = []string{
				logo,
				cart,
			}
		}
	default:
		if m.checkout {
			tabs = []string{
				back,
				logo,
				cart,
			}
		} else {
			tabs = []string{
				logo,
				shop,
				about,
				faq,
				cart,
			}
		}
	}

	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(m.renderer.NewStyle().Foreground(m.theme.Border())).
		Row(tabs...).
		Width(m.widthContainer).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.theme.Base().
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center)
		}).
		Render()
}
