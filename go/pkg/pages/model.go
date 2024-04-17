package pages

import (
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
		pages: []Page{
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
    Update(m Model, raw tea.Msg) (bool, tea.Model, tea.Cmd)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) systemUpdates(raw tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg := raw.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

        if m.width < MIN_WIDTH || m.height < MIN_HEIGHT {
            m.currentPage = 0
        } else {
            // TODO: we need to implement history
            m.currentPage = 1
        }

		return true, m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return true, m, tea.Quit
        }
    }
    return false, m, nil
}

func (m Model) Update(raw tea.Msg) (tea.Model, tea.Cmd) {

    if handled, model, cmd := m.systemUpdates(raw); handled {
        return model, cmd
    }

	// Not sure this is great... but it's kind of nice to all be in the same place
	page := m.pages[m.currentPage]
    if handled, model, cmd := page.Update(m, raw); handled {
        return model, cmd
	}

	return m, nil
}

func (m Model) View() string {
	page := m.pages[m.currentPage]

	theme := GetTheme(m.renderer)
    title := theme.ActivePage().Render(page.Title())

	pageStyle := m.renderer.NewStyle()

	return lipgloss.JoinVertical(
        lipgloss.Top,
        title,
        pageStyle.Render(page.Render(&m)))
}

// func (m *Model) SetProductCount
