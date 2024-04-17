package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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
    goToShipping = 2
    goToCC = 3
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

    log.Warn("initial state", "state", state, "email", goToEmail, "shippingState", goToShipping, "cc", goToCC)
    if state == goToEmail {
        model.order.count = 1
        model.order.widget = api.GetWidgets()[0]
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
		} else if m.currentPage == 0 {
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
		case "ctrl+c":
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

func (m Model) createTitle() string {
	theme := GetTheme(m.renderer)

    titles := make([]string, 0)
    for i, page := range m.pages {
        if i == 0 {
            continue
        }
        current := i == m.currentPage
        style := theme.Page()
        if current {
            style = theme.ActivePage()
        }

        title := style.MarginLeft(1).Render(fmt.Sprintf("> %s", page.Title()))

        titles = append(titles, title)
    }

    return lipgloss.JoinHorizontal(0, titles...)
}

func (m Model) View() string {
	page := m.pages[m.currentPage]
    log.Warn("I am on page", "title", page.Title(), "currentPage", m.currentPage)

	pageStyle := m.renderer.NewStyle()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.createTitle(),
		pageStyle.Render(page.Render(&m)))
}

// func (m *Model) SetProductCount
