package components

import (
	"strings"
)

type Layout struct {
	Title      string
	Content    string
	Footer     string
	Width      int
	Height     int
	ShowBorder bool
}

func NewLayout(title string) *Layout {
	return &Layout{
		Title:      title,
		ShowBorder: true,
	}
}

func (l *Layout) Render() string {
	var b strings.Builder

	// Header with title
	if l.Title != "" {
		header := TitleStyle.Render(l.Title)
		b.WriteString(header)
		b.WriteString("\n\n")
	}

	// Main content
	if l.Content != "" {
		if l.ShowBorder {
			content := BoxStyle.Width(l.Width - 6).Render(l.Content)
			b.WriteString(content)
		} else {
			b.WriteString(l.Content)
		}
		b.WriteString("\n")
	}

	// Footer
	if l.Footer != "" {
		footer := FooterStyle.Width(l.Width - 4).Render(l.Footer)
		b.WriteString(footer)
	}

	// Apply base style and dimensions
	finalStyle := BaseStyle
	if l.Width > 0 {
		finalStyle = finalStyle.Width(l.Width)
	}
	if l.Height > 0 {
		finalStyle = finalStyle.Height(l.Height)
	}

	return finalStyle.Render(b.String())
}
