package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/terminalhq/terminal/go/pkg/api"
)

const (
    BREAD_CRUMB_HEIGHT = 3
    HELP_MENU = 3

    MIN_WIDTH_NOT_MET_PAGE = 0
    PRODUCT_PAGE = 1
    EMAIL_PAGE = 2
    SHIPPING_PAGE = 3
    CC_PAGE = 4
    SUMMARY_PAGE = 5
)



// type currentPage int
//
// const (
// 	Cart currentPage = iota
// 	Product
// )

type OrderInfo struct {
	count  int
	product api.Product
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

	shippingState   ShippingState
	creditCardState CreditCardState

    Dialog *string
}

var defaultShippingState = ShippingState{
	Name:      "default Name",
	AddrLine1: "default AddrLine1",
	AddrLine2: "default AddrLine2",
	City:      "default City",
	State:     "default State",
	Zip:       "default Zip",
}

var defaultCreditCardState = CreditCardState{
	Name: "default Name",

	CC:       "default CC",
	CVC:      "default CVC",
	ExpMonth: "default ExpMonth",
	ExpYear:  "default ExpYear",
}

var defaultEmail = "piq@called.it"

const (
	goToEmail    = 1
	goToShipping = 2
	goToCC       = 3
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

	model := &Model{
		renderer:    renderer,
		currentPage: PRODUCT_PAGE,
		theme:       GetTheme(renderer),
		width:       0,
		email:       "",
		height:      0,
        Dialog:      nil,
		pages: []Page{
			&MinWidthPage{},
			&ProductPage{},
			NewEmailPage(),
			NewShippingPage(),
			NewCreditCardPage(),
		},
		order: OrderInfo{
			count:  0,
			product: api.GetProducts()[0],
		},
	}

	state := stateToNumber(toState)

	log.Warn("initial state", "state", state, "email", goToEmail, "shippingState", goToShipping, "cc", goToCC)
	if state == goToEmail {
		model.order.count = 1
		model.order.product = api.GetProducts()[0]
		model.currentPage = EMAIL_PAGE
	}
	if state == goToShipping {
		model.email = defaultEmail
		model.currentPage = SHIPPING_PAGE
	}

	if state == goToCC {
		model.email = defaultEmail
		model.shippingState = defaultShippingState
		model.currentPage = CC_PAGE
	}

	return model
}

func (m *Model) GetMaxPageHeight() int {
    return m.height - (BREAD_CRUMB_HEIGHT + HELP_MENU)
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
			m.currentPage = MIN_WIDTH_NOT_MET_PAGE
		} else if m.currentPage == MIN_WIDTH_NOT_MET_PAGE {
            // PLEASE CHANGE THIS
			m.currentPage = PRODUCT_PAGE
		}

		return true, m, nil
	case NavigateProduct:
		m.currentPage = PRODUCT_PAGE
		return true, m, nil
	case NavigateEmail:
		m.currentPage = EMAIL_PAGE
		return true, m, nil
	case NavigateShipping:
		m.currentPage = SHIPPING_PAGE
		return true, m, nil
	case Dialog:
		m.Dialog = &msg.msg
		return true, m, nil
	case NavigateCC:
		m.currentPage = CC_PAGE
		return true, m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
            m.Dialog = nil
            return true, m, nil
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

    titleContainer := lipgloss.NewStyle().
        Margin(1, 0, 1, 0)

    if m.currentPage == 0 {
        return titleContainer.Render(" ")
    }

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

	return titleContainer.Render(lipgloss.JoinHorizontal(0, titles...))
}

func (m Model) View() string {

    var renderedPage string
    if m.Dialog != nil && len(*m.Dialog) > 0 {
        renderedPage = DisplayDialog(m, *m.Dialog)
    } else {
        page := m.pages[m.currentPage]
        log.Warn("I am on page", "title", page.Title(), "currentPage", m.currentPage)

        pageStyle := m.renderer.NewStyle()
        renderedPage = pageStyle.Render(page.Render(&m))
    }

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.createTitle(),
        renderedPage,
        helpMenu(m),
    )
}

