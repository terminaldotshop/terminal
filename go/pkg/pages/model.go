package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminalhq/terminal/go/pkg/api"
)

// type currentPage int
//
// const (
// 	Cart currentPage = iota
// 	Widget
// )

type OrderInfo struct {
	count  int
	widget api.Widget
}

type Model struct {
	currentPage int
	pages       []Page

	width  int
	height int

	renderer *lipgloss.Renderer
	theme    Theme
	order    OrderInfo
}

func NewModel() *Model {
	renderer := lipgloss.DefaultRenderer()

	return &Model{
		renderer:    renderer,
		currentPage: 1,
		theme:       GetTheme(renderer),
		width:       0,
		height:      0,
		pages:       []Page{
            &MinWidthPage{},
            &WidgetPage{},
        },
		order: OrderInfo{
			count:  0,
			widget: api.GetWidgets()[0],
		},
	}
}

type Page interface {
	Title() string
	Render(m *Model) string
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
