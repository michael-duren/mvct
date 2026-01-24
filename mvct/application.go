package mvct

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Application is the main framework orchestrator
// M is the global model type
type Application[M any] struct {
	router         *Router
	settings       map[string]any
	width          int
	height         int
	globalHandlers []GlobalHandler
	model          M
}

// Config holds application configuration
type Config struct {
	DefaultRoute string
	Settings     map[string]any
}

func NewApplication[M any](config Config, model M) *Application[M] {
	app := &Application[M]{
		router:         NewRouter(),
		globalHandlers: []GlobalHandler{},
		model:          model,
	}

	if config.DefaultRoute != "" {
		app.router.SetDefault(config.DefaultRoute)
	}

	return app
}

// RegisterController adds a controller for a route
func (a *Application[M]) RegisterController(path string, controller Controller[any]) {
	a.router.Register(path, controller)
}

// Use adds middleware
func (a *Application[M]) Use(m Middleware) {
	a.router.Use(m)
}

// UseGlobalHandler adds a handler that runs before routing
func (a *Application[M]) UseGlobalHandler(handler GlobalHandler) {
	a.globalHandlers = append(a.globalHandlers, handler)
}

// Model returns the global model
func (a *Application[M]) Model() any {
	return a.model
}

// SetModel updates the global model
func (a *Application[M]) SetModel(model M) {
	a.model = model
}

// Width returns the current window width
func (a *Application[M]) Width() int {
	return a.width
}

// Height returns the current window height
func (a *Application[M]) Height() int {
	return a.height
}

// ErrSettingNotFound is returned when a setting is not found
type ErrSettingNotFound struct {
	key string
}

func (e ErrSettingNotFound) Error() string {
	return fmt.Sprintf("setting with key: %s was not found", e.key)
}

// GetSetting retrieves a setting by key
func (a *Application[M]) GetSetting(s string) (any, error) {
	value, ok := a.settings[s]
	if !ok {
		return nil, &ErrSettingNotFound{key: s}
	}
	return value, nil
}

// Init implements tea.Model
func (a *Application[M]) Init() tea.Cmd {
	return a.router.Current().Init()
}

// Update implements tea.Model
func (a *Application[M]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		// Run global handlers
		for _, handler := range a.globalHandlers {
			if cmd := handler.Handle(msg); cmd != nil {
				return a, cmd
			}
		}

	case NavigateMsg:
		cmd, err := a.router.Navigate(msg.Route)
		if err != nil {
			// Could handle error here
			return a, nil
		}
		return a, cmd
	}

	// Route to current controller
	cmd := a.router.Current().Update(msg)
	return a, cmd
}

// View implements tea.Model
func (a *Application[M]) View() string {
	return a.router.Current().View()
}

func (a *Application[M]) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}
