package mvct

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type KeyMsgHandler func(msg KeyMsg) Cmd

// Controller is the interface all controllers must implement
type Controller interface {
	// Init is called when the controller is first activated
	// the caller can register the key handlers they want to be used
	// while on this route
	Init(handlers KeyHandlers) Cmd

	// View renders the controller's view
	View() string
}

// Msg wraps the tea.Msg with additional context
type Msg struct {
	Inner   tea.Msg
	Context context.Context
}

// Cmd is a command that returns a Msg
type Cmd func() Msg

// Navigate should be called by controllers to change routes
func Navigate(route string) Cmd {
	return func() Msg {
		return Msg{
			Inner: NavigateMsg{Route: route},
		}
	}
}
