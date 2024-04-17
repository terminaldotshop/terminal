package pages

import (
	"strings"

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

	// Customer information
	order OrderInfo
	email string

    shippingState ShippingState
    creditCardState CreditCardState
}

var defaultShippingState = ShippingState{
	Name: "default Name",
	AddrLine1: "default AddrLine1",
	AddrLine2: "default AddrLine2",
	City: "default City",
	State: "default State",
	Zip: "default Zip",
}

var defaultCreditCardState = CreditCardState{
	Name: "default Name",

	CC: "default CC",
	CVC: "default CVC",
	ExpMonth: "default ExpMonth",
	ExpYear: "default ExpYear",
}

var defaultEmail = "piq@called.it"

const (
    goToEmail = 1
    goToShipping
    goToCC
)

func stateToNumber(toState string) int {
    switch strings.ToLower(toState) {
    case "email":
        return 1
    case "shipping":
        return 2
    case "cc":
        return 3
    }
    return 0
}

func NewModel(toState string) *Model {
	renderer := lipgloss.DefaultRenderer()

    model :=  &Model{
		renderer:    renderer,
		currentPage: 1,
		theme:       GetTheme(renderer),
		width:       0,
		email:       "",
		height:      0,
		pages: []Page{
			&MinWidthPage{},
			&WidgetPage{},
			NewEmailPage(),
			NewShippingPage(),
			NewCreditCardPage(),
		},
		order: OrderInfo{
			count:  0,
			widget: api.GetWidgets()[0],
		},
	}

    state := stateToNumber(toState)
    if state == goToEmail {
        model.currentPage = 2
    }
    if state == goToShipping {
        model.email = defaultEmail
        model.currentPage = 3
    }

    if state == goToCC {
        model.email = defaultEmail
        model.shippingState = defaultShippingState
        model.currentPage = 4
    }

    return model
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
	case BeginCheckout:
		m.currentPage = 2
		return true, m, nil
	case StartShipping:
		m.currentPage = 3
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
