package tui

func (m model) BreadcrumbsView() string {
	accent := m.theme.TextAccent().Render
	base := m.theme.Base().Render
	sep := m.theme.Base().Render("/")

	var labels []string
	switch m.size {
	case small:
		fallthrough
	case medium:
		labels = []string{"cart", "pay", "ship", "confirm"}
	default:
		labels = []string{"cart", "payment", "shipping", "confirmation"}
	}

	var selected int
	switch m.page {
	case cartPage:
		selected = 0
	case paymentPage:
		selected = 1
	case shippingPage:
		selected = 2
	case confirmPage:
		selected = 3
	default:
		return ""
	}

	items := []string{}
	for i, label := range labels {
		if i == selected {
			items = append(items, accent(label))
			items = append(items, sep)
		} else {
			items = append(items, base(label))
			items = append(items, sep)
		}
	}

	// remove last separator
	items = items[:len(items)-1]

	return m.theme.Base().
		MarginTop(1).
		MarginBottom(1).
		PaddingLeft(1).
		Render(items...)
}
