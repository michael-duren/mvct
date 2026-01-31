package mvct

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type KeyMsgHandler func(msg KeyMsg) Cmd

// Controller is the base interface all controllers must implement
// M is the model type associated with the controller
// that should be passed to the View method
type Controller interface {
	// Init is called when the controller is first activated
	Init() Cmd

	// View renders the controller's view
	View() string

	// Returns the model associated with the controller
	GetModel() any

	// Registers key handlers specific to this controller
	RegisterKeyHandlers(map[string]KeyMsgHandler)
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

type Result[T any] struct {
	Success bool
	Value   T
}

func NewResult[T any](val T, success bool) *Result[T] {
	return &Result[T]{
		Value:   val,
		Success: success,
	}
}
