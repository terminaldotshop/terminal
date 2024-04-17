package pages

import "github.com/charmbracelet/lipgloss"

type Theme interface {
	PriceForeground() lipgloss.Style
	DescForeground() lipgloss.Style
	TitleForeground() lipgloss.Style
	NormalForeground() lipgloss.Style
	ActivePriceForeground() lipgloss.Style
	ActiveDescForeground() lipgloss.Style
	ActiveTitleForeground() lipgloss.Style
	HelpSpecialForeground() lipgloss.Style
	HRB() lipgloss.Style
	ActiveNormalForeground() lipgloss.Style

	Page() lipgloss.Style
	ActivePage() lipgloss.Style
}

type BasicTheme struct {
	price        lipgloss.Style
	desc         lipgloss.Style
	help         lipgloss.Style
	title        lipgloss.Style
	normal       lipgloss.Style
	page         lipgloss.Style
	hrb          lipgloss.Style
	activePrice  lipgloss.Style
	activeDesc   lipgloss.Style
	activeTitle  lipgloss.Style
	activeNormal lipgloss.Style
	activePage   lipgloss.Style
}

func (b *BasicTheme) HRB() lipgloss.Style {
	return b.hrb.Copy()
}

func (b *BasicTheme) HelpSpecialForeground() lipgloss.Style {
	return b.help.Copy()
}

func (b *BasicTheme) PriceForeground() lipgloss.Style {
	return b.activeDesc.Copy()
}

func (b *BasicTheme) DescForeground() lipgloss.Style {
	return b.desc.Copy()
}

func (b *BasicTheme) TitleForeground() lipgloss.Style {
	return b.title.Copy()
}

func (b *BasicTheme) NormalForeground() lipgloss.Style {
	return b.normal.Copy()
}

func (b *BasicTheme) ActivePriceForeground() lipgloss.Style {
	return b.activeDesc.Copy()
}

func (b *BasicTheme) ActiveDescForeground() lipgloss.Style {
	return b.activeDesc.Copy()
}

func (b *BasicTheme) ActiveTitleForeground() lipgloss.Style {
	return b.activeTitle.Copy()
}

func (b *BasicTheme) ActiveNormalForeground() lipgloss.Style {
	return b.activeNormal.Copy()
}

func (b *BasicTheme) Page() lipgloss.Style {
	return b.page.Copy()
}

func (b *BasicTheme) ActivePage() lipgloss.Style {
	return b.activePage.Copy()
}

func GetTheme(renderer *lipgloss.Renderer) Theme {
	return &BasicTheme{
		desc:   renderer.NewStyle(),
		title:  renderer.NewStyle(),
		normal: renderer.NewStyle(),
		page: renderer.NewStyle().
			Foreground(lipgloss.Color("#4e545d")).
			MarginLeft(4),
		activeDesc: renderer.NewStyle(),
		hrb: renderer.NewStyle().
			Bold(true).
            Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("#EC465A")),
		activeTitle: renderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#b294bb")).
			Underline(true),
		activeNormal: renderer.NewStyle(),
		activePage: renderer.NewStyle().
			Foreground(lipgloss.Color("#99cc99")).
			MarginLeft(2).
			Bold(true),
		help: renderer.NewStyle().
			Foreground(lipgloss.Color("#FF5C00")).
			MarginLeft(2).
			Bold(true),
	}
}
