package mvct

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

type KeyHandlers = map[string]KeyMsgHandler

// Application is the main framework orchestrator
// M is the global model type
type Application[M any] struct {
	router         *Router
	settings       map[string]any
	width          int
	height         int
	globalHandlers []GlobalHandler
	// key handlers are registered during the init method when the application
	// starts or when a route changes
	keyHandlers KeyHandlers
	// message handlers are registered by type - auto-discovered via reflection
	msgHandlers map[reflect.Type]reflect.Value
	model       M
	layoutFunc  func(content string, width, height int) string

	Errors []error
}

// Config holds application configuration
type Config struct {
	DefaultRoute string
	Settings     map[string]any
}

func NewApplication[M any](config Config, model M) *Application[M] {
	slog.Info("Creating new application", "default_route", config.DefaultRoute)
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
func (a *Application[M]) RegisterController(path string, controller Controller) {
	slog.Debug("Registering controller", "path", path, "controller_type", reflect.TypeOf(controller))
	a.router.Register(path, controller)
}

// Use adds middleware
func (a *Application[M]) Use(m Middleware) {
	slog.Debug("Registering middleware", "middleware_type", reflect.TypeOf(m))
	a.router.Use(m)
}

// UseGlobalHandler adds a handler that runs before routing
func (a *Application[M]) UseGlobalHandler(handler GlobalHandler) {
	slog.Debug("Registering global handler", "handler_type", reflect.TypeOf(handler))
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
	slog.Debug("Initializing application")
	cmd := a.router.Current().Init(a.keyHandlers)
	a.scanMessageHandlers()
	return unwrapCmd(cmd)
}

func (a *Application[M]) View() string {
	slog.Debug("Application View called")
	content := a.router.Current().View()

	if a.layoutFunc != nil {
		return a.layoutFunc(content, a.width, a.height)
	}

	return content
}

func (a *Application[M]) SetLayout(fn func(content string, width, height int) string) {
	slog.Debug("Setting application layout")
	a.layoutFunc = fn
}

func (a *Application[M]) Run() error {
	slog.Info("Starting application run loop")
	p := tea.NewProgram(a)
	_, err := p.Run()
	slog.Info("Application stopped")
	return err
}

func (a *Application[M]) logHandlers() {
	slog.Debug("=== Registered Handlers ===")

	for key, handler := range a.keyHandlers {
		fnName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		slog.Debug("key handler",
			"key", key,
			"function", fnName,
		)
	}

	for msgType, method := range a.msgHandlers {
		fnName := runtime.FuncForPC(method.Pointer()).Name()
		slog.Debug("message handler",
			"handles", msgType.String(),
			"method", fnName,
		)
	}
}
