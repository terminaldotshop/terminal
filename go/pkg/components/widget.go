package components

import "github.com/terminalhq/terminal/go/pkg/window"

type Widget struct {
	Name        string
	Description string

	AsciiArtBig   string
	AsciiArtSmall string
}

type InlineWidget struct {
    Widget
}

func (w *InlineWidget) Render() (window.Location, [][]*window.Cell) {
}


