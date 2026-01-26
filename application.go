package mvct

import (
	"fmt"
	"reflect"
	"strings"

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
	// key handlers are registered via the RegisterKeyHandlers method and are mapped by key string, use convience methods to create
	// these maps
	keyHandlers map[string]KeyMsgHandler
	// message handlers are registered by type - auto-discovered via reflection
	msgHandlers map[reflect.Type]reflect.Value
	model       M
	Errors      []error
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
func (a *Application[M]) RegisterController(path string, controller Controller) {
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
	cmd := a.router.Current().Init()
	return unwrapCmd(cmd)
}

// Update implements tea.Model
func (a *Application[M]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case tea.KeyMsg:
		for _, handler := range a.globalHandlers {
			if cmd := handler.Handle(msg); cmd != nil {
				return a, cmd
			}
		}
	}

	wrappedMsg := wrapMsg(msg)

	if navMsg, ok := wrappedMsg.Inner.(NavigateMsg); ok {
		cmd, err := a.router.Navigate(navMsg.Route)
		if err != nil {
			a.Errors = append(a.Errors, err)
			return a, nil
		}

		a.scanMessageHandlers()
		return a, cmd
	}

	// local key msgs
	if keyMsg, ok := wrappedMsg.Inner.(KeyMsg); ok {
		if handler, exists := a.keyHandlers[keyMsg.Type.String()]; exists {
			return a, unwrapCmd(handler(keyMsg))
		}
	}

	msgType := reflect.TypeOf(wrappedMsg.Inner)
	if handler, exists := a.msgHandlers[msgType]; exists {
		results := handler.Call([]reflect.Value{reflect.ValueOf(wrappedMsg.Inner)})
		if len(results) > 0 && !results[0].IsNil() {
			return a, unwrapCmd(results[0].Interface().(Cmd))
		}
	}

	// Route to current controller via internalUpdate
	currentController := a.router.Current()
	if bc, ok := currentController.(interface{ internalUpdate(Msg) Cmd }); ok {
		cmd := bc.internalUpdate(wrappedMsg)
		return a, unwrapCmd(cmd)
	}

	return a, nil
}

// gets the appropriate handler based on the paramter
func (a *Application[M]) scanMessageHandlers() {
	a.keyHandlers = make(map[string]KeyMsgHandler)
	a.msgHandlers = make(map[reflect.Type]reflect.Value)

	ctlr := a.router.Current()
	if ctlr == nil {
		return
	}

	val := reflect.ValueOf(ctlr)
	typ := val.Type()

	for i := 0; i < val.NumMethod(); i++ {
		method := val.Method(i)
		methodName := typ.Method(i).Name
		methodType := method.Type()

		if strings.HasPrefix(methodName, "OnKey") {
			continue
		}

		if strings.HasPrefix(methodName, "On") {
			if methodType.NumIn() == 1 && methodType.NumOut() == 1 {
				msgType := methodType.In(0)
				a.msgHandlers[msgType] = method
			}
		}
	}
	ctlr.RegisterKeyHandlers(a.keyHandlers)
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
