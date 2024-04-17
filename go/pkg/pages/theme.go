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
	ActiveNormalForeground() lipgloss.Style

	Page() lipgloss.Style
	ActivePage() lipgloss.Style
}

type BasicTheme struct {
	price        lipgloss.Style
	desc         lipgloss.Style
	title        lipgloss.Style
	normal       lipgloss.Style
	page         lipgloss.Style
	activePrice  lipgloss.Style
	activeDesc   lipgloss.Style
	activeTitle  lipgloss.Style
	activeNormal lipgloss.Style
	activePage   lipgloss.Style
}

func (b *BasicTheme) PriceForeground() lipgloss.Style {
	return b.activeDesc
}

func (b *BasicTheme) DescForeground() lipgloss.Style {
	return b.desc
}

func (b *BasicTheme) TitleForeground() lipgloss.Style {
	return b.title
}

func (b *BasicTheme) NormalForeground() lipgloss.Style {
	return b.normal
}

func (b *BasicTheme) ActivePriceForeground() lipgloss.Style {
	return b.activeDesc
}

func (b *BasicTheme) ActiveDescForeground() lipgloss.Style {
	return b.activeDesc
}

func (b *BasicTheme) ActiveTitleForeground() lipgloss.Style {
	return b.activeTitle
}

func (b *BasicTheme) ActiveNormalForeground() lipgloss.Style {
	return b.activeNormal
}

func (b *BasicTheme) Page() lipgloss.Style {
	return b.page
}

func (b *BasicTheme) ActivePage() lipgloss.Style {
	return b.activePage
}

func GetTheme(renderer *lipgloss.Renderer) Theme {
	return &BasicTheme{
		desc:   renderer.NewStyle(),
		title:  renderer.NewStyle(),
		normal: renderer.NewStyle(),
		page: renderer.NewStyle().
			Foreground(lipgloss.Color("#4e545d")).
			MarginLeft(4),
		activeDesc:   renderer.NewStyle(),
		activeTitle:  renderer.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color("#b294bb")).
            Underline(true),
		activeNormal: renderer.NewStyle(),
		activePage: renderer.NewStyle().
			Foreground(lipgloss.Color("#99cc99")).
			MarginLeft(2).
			Bold(true),
	}
}
