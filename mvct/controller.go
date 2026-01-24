package mvct

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Controller is the base interface all controllers must implement
// M is the model type associated with the controller
// that should be passed to the View method
type Controller[M any] interface {
	// Init is called when the controller is first activated
	Init() tea.Cmd

	// Update handles messages for this controller
	Update(msg tea.Msg) tea.Cmd

	// View renders the controller's view
	View() string

	// Returns the model associated with the controller
	GetModel() M
}

// BaseController provides common functionality
// and global access to the Application instance
type BaseController[M any] struct {
	app *Application[M]
}

func (bc *BaseController[M]) SetApp(app *Application[M]) {
	bc.app = app
}

func (bc *BaseController[M]) App() *Application[M] {
	return bc.app
}

// Navigate changes the current route
func (bc *BaseController[M]) Navigate(route string) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{Route: route}
	}
}

// NavigateMsg signals a route change
type NavigateMsg struct {
	Route string
}
