package pages

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/terminalhq/terminal/go/pkg/api"
)

const (
	BREAD_CRUMB_HEIGHT = 3
	HELP_MENU          = 3

	MIN_WIDTH_NOT_MET_PAGE = 0
	PRODUCT_PAGE           = 1
	EMAIL_PAGE             = 2
	SHIPPING_PAGE          = 3
	CC_PAGE                = 4
	CC_ADDR_PAGE           = 5
	CONFIRM_PAGE           = 6
	ANIMATION_PAGE         = 7
)

type OrderInfo struct {
	count   int
	product *api.Product
}

type Model struct {
	testing bool

	userToken   string
	currentPage int
	pages       []Page

	minWidth bool

	width  int
	height int

	renderer *lipgloss.Renderer
	theme    Theme

	// Customer information
	order OrderInfo
	email string

	creditCard      api.CreditCard
	stripeCardToken *api.StripeCardToken

	differentBillingAddress bool

	shippingAddress api.Address
	billingAddress  api.Address

	Dialog *string
}

func NewModel(
	renderer *lipgloss.Renderer,
	width int,
	height int,
	publicKey string,
) (Model, error) {
	log.Warn("Creating new model")
	product, err := api.FetchOneProduct()
	if err != nil {
		return Model{}, err
	}

	log.Warn("Creating New Product Page")
	productPage := NewProductPage(product)

	log.Warn("Fetching User Token")
	userToken, err := api.FetchUserToken(publicKey)
	if err != nil {
		return Model{}, err
	}

	log.Warn("starting terminal.shop", "page", PRODUCT_PAGE, "title", productPage.Title())
	model := Model{
		testing:     false,
		userToken:   userToken.AccessToken,
		width:       width,
		height:      height,
		renderer:    renderer,
		currentPage: PRODUCT_PAGE,
		theme:       GetTheme(renderer),
		email:       "",
		Dialog:      nil,
		pages: []Page{
			NewMinWidthPage(),
			productPage,
			NewEmailPage(),
			NewShippingPage(),
			NewCreditCardPage(),
			NewCreditCardAddress(),
			NewConfirmPage(),
			NewAnimationPage(),
		},
		order: OrderInfo{
			count:   0,
			product: productPage.Product,
		},
	}

	model.pages[model.currentPage].Enter(model)
	log.Warn("starting terminal.shop", "page", model.currentPage, "title", model.pages[model.currentPage].Title())

	return model, nil
}

func NewLocalModel(toState string) *Model {
	renderer := lipgloss.DefaultRenderer()

	product, err := api.FetchOneProduct()
	if err != nil {
		panic("no")
	}

	productPage := NewProductPage(product)

	model := &Model{
		testing:     true,
		renderer:    renderer,
		currentPage: PRODUCT_PAGE,
		theme:       GetTheme(renderer),
		width:       0,
		email:       "",
		height:      0,
		Dialog:      nil,
		pages: []Page{
			&MinWidthPage{},
			productPage,
			NewEmailPage(),
			NewShippingPage(),
			NewCreditCardPage(),
			NewCreditCardAddress(),
			NewConfirmPage(),
			// TODO: Add a page to show that order worked. Animate the coffee
			NewAnimationPage(),
		},
		order: OrderInfo{
			count:   0,
			product: productPage.Product,
		},
	}

	state := stateToNumber(toState)

	log.Warn("initial state",
		"state", state,
		"email", goToEmail,
		"shippingState", goToShipping,
		"cc", goToCC,
		"cc-addr", goToCCAddr,
		"confirm", goToConfirm,
	)

	log.Warn("test", "state", state, "email", goToEmail, "shipping", goToShipping, "cc", goToCC, "ccaddr", goToCCAddr, "con", goToConfirm)
	if state >= goToEmail {
		product, err := api.FetchOneProduct()
		if err != nil {
			log.Fatal("Failed to fetch. Tag @thdxr on x.com")
		}

		log.Warn("order info")
		model.order.count = 1
		model.order.product = product
		model.currentPage = EMAIL_PAGE
	}

	if state >= goToShipping {
		log.Warn("email")
		model.email = defaultEmail
		model.currentPage = SHIPPING_PAGE
	}

	if state >= goToCC {
		log.Warn("shipping", "cc page?", CC_PAGE)
		model.shippingAddress = defaultShippingState
		model.currentPage = CC_PAGE
	}

	if state >= goToCCAddr {
		log.Warn("cc")
		model.creditCard = defaultCreditCard
		model.currentPage = CC_ADDR_PAGE
	}

	if state >= goToConfirm {
		log.Warn("cc addr")
		model.billingAddress = defaultBillingAddress
		model.currentPage = CONFIRM_PAGE
	}

	if state >= goToAnimation {
		model.currentPage = ANIMATION_PAGE
	}

	model.pages[model.currentPage].Enter(*model)
	log.Warn("starting terminal.shop", "page", model.currentPage, "title", model.pages[model.currentPage].Title())

	return model
}

func (m *Model) GetMaxPageHeight() int {
	return m.height - (BREAD_CRUMB_HEIGHT + HELP_MENU)
}

func (m *Model) GetMaxPageWidth() int {
	return m.width
}

type Page interface {
	Enter(m Model)
	Exit(m Model) Model
	Title() string
	Render(m *Model) string
	Update(m Model, raw tea.Msg) (bool, tea.Model, tea.Cmd)
}

func (m Model) Init() tea.Cmd {
	api.Init(m.testing)

	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return NextFrameMsg{}
	})

}

func nav(m Model, newPage int) Model {
	log.Warn("navigation event", "from", m.currentPage, "to", newPage)
	m = m.pages[m.currentPage].Exit(m)
	m.currentPage = newPage
	m.pages[m.currentPage].Enter(m)
	return m
}

func (m Model) systemUpdates(raw tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg := raw.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.minWidth = m.width < MIN_WIDTH || m.height < MIN_HEIGHT
		return true, m, nil

	case NavigateProduct:
		return true, nav(m, PRODUCT_PAGE), nil
	case NavigateEmail:
		return true, nav(m, EMAIL_PAGE), nil
	case NavigateShipping:
		return true, nav(m, SHIPPING_PAGE), nil
	case NavigateCC:
		return true, nav(m, CC_PAGE), nil
	case NavigateCCAddress:
		return true, nav(m, CC_ADDR_PAGE), nil
	case NavigateConfirm:
		return true, nav(m, CONFIRM_PAGE), nil
	case NavigateAnimation:
		return true, nav(m, ANIMATION_PAGE), tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return NextFrameMsg{}
		})

	case Dialog:
		m.Dialog = &msg.msg
		return true, m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			if m.currentPage >= EMAIL_PAGE {
				m.currentPage -= 1
				if m.currentPage == CC_ADDR_PAGE && !m.differentBillingAddress {
					m.currentPage -= 1
				}
			}
			m.pages[m.currentPage].Enter(m)
			return true, m, nil
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

	titleContainer := m.renderer.NewStyle().
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
		if m.minWidth {
			page = m.pages[0]
		}

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
