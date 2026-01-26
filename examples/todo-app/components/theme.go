package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor = lipgloss.Color("#EC4899") // Pink
	AccentColor    = lipgloss.Color("#10B981") // Green
	TextColor      = lipgloss.Color("#F3F4F6")
	MutedColor     = lipgloss.Color("#9CA3AF")
	BgColor        = lipgloss.Color("#1F2937")
	BorderColor    = lipgloss.Color("#374151")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BgColor)

	// Header style
	HeaderStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor)

	// List item styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true).
				PaddingLeft(1)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				PaddingLeft(1)

	// Input style
	InputStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginTop(1).
			Padding(0, 1)

	// Footer style
	FooterStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(BorderColor).
			MarginTop(1).
			Padding(1, 2)

	// Container/Box style
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			MarginTop(1)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)
)
