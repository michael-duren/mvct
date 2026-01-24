package mvc

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
type BaseController struct {
	app *Application
}

func (bc *BaseController) SetApp(app *Application) {
	bc.app = app
}

func (bc *BaseController) App() *Application {
	return bc.app
}

// Navigate changes the current route
func (bc *BaseController) Navigate(route string) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{Route: route}
	}
}

// NavigateMsg signals a route change
type NavigateMsg struct {
	Route string
}
