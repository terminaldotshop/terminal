package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"test.terminal.shop/pkg/shop"
)

// type currentPage int
//
// const (
// 	Cart currentPage = iota
// 	Widget
// )

type WidgetOrderInfo struct {
	count  int
	widget shop.Widget
}

type CartInfo struct {
	totalItems int
	widgets    []WidgetOrderInfo
}

type Model struct {
	currentPage int
	pages       []Page

	renderer *lipgloss.Renderer
	theme    Theme
	cart     CartInfo
	width    int
	height   int
}

func NewModel() *Model {
	widgets := shop.GetWidgets()

	renderer := lipgloss.DefaultRenderer()

	return &Model{
		renderer:    renderer,
		currentPage: 0,
		theme:       GetTheme(renderer),
		cart: CartInfo{
			totalItems: 4,
			widgets: []WidgetOrderInfo{
				{count: 3, widget: shop.GetWidgets()[1]},
				{count: 1, widget: shop.GetWidgets()[0]},
			},
		},

		pages: []Page{
			&CartPage{},
			&WidgetPage{
				widget: &widgets[0],
			},
		},
	}
}

type Page interface {
	Title() string
	Render(m *Model) string
	// Update(m *Model, msg tea.Msg) string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(raw tea.Msg) (tea.Model, tea.Cmd) {
	// Not sure this is great... but it's kind of nice to all be in the same place
	page := m.pages[m.currentPage]
	switch page := page.(type) {
	case *WidgetPage:
		_ = page

		switch msg := raw.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "left":
				// This would be where we can change the amount of this page
			}
		}
	}

	switch msg := raw.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.currentPage = (m.currentPage + 1) % len(m.pages)
		}
	}
	return m, nil
}

func (m Model) View() string {
	page := m.pages[m.currentPage]

	theme := GetTheme(m.renderer)

	titles := []string{}
	for idx, page := range m.pages {
		if idx == m.currentPage {
			titles = append(titles, theme.ActivePage().Render(
				fmt.Sprintf("* %s", page.Title())),
			)
		} else {
			titles = append(titles, theme.Page().Render(page.Title()))
		}
	}

	headers := lipgloss.JoinHorizontal(lipgloss.Left, titles...)
	pageStyle := m.renderer.NewStyle()

	return fmt.Sprintf("%s\n%s", headers, pageStyle.Render(page.Render(&m)))
}

// func (m *Model) SetProductCount
