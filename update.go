package mvct

import (
	"log/slog"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model
func (a *Application[M]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("Application Update", "msg_type", reflect.TypeOf(msg), "msg", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return a.handleWindowResize(msg)
	}

	wrappedMsg := wrapMsg(msg)

	switch inner := wrappedMsg.Inner.(type) {
	case NavigateMsg:
		return a.handleNavigate(inner)
	case KeyMsg:
		if cmd, ok := a.handleKeyMsg(inner, wrappedMsg); ok {
			return a, cmd
		}
	default:
		if cmd, ok := a.handleControllerMsg(wrappedMsg); ok {
			return a, cmd
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return a.handleGlobalKeyMsg(keyMsg)
	}

	slog.Debug("No handler found for message", "msg_type", reflect.TypeOf(wrappedMsg.Inner).String())
	return a, nil
}

func (a *Application[M]) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	slog.Debug("Window resize", "width", msg.Width, "height", msg.Height)
	a.width = msg.Width
	a.height = msg.Height
	return a, nil
}

func (a *Application[M]) handleNavigate(msg NavigateMsg) (tea.Model, tea.Cmd) {
	slog.Info("Processing navigation", "route", msg.Route)
	cmd, err := a.router.Navigate(a.keyHandlers, msg.Route)
	if err != nil {
		slog.Error("Navigation failed", "error", err)
		a.Errors = append(a.Errors, err)
		return a, nil
	}
	a.scanMessageHandlers()
	slog.Debug("Navigation successful", "new_route", msg.Route)
	a.logHandlers()
	return a, cmd
}

func (a *Application[M]) handleKeyMsg(msg KeyMsg, wrappedMsg Msg) (tea.Cmd, bool) {
	var cmds []tea.Cmd

	// 1. Generic OnKeyMsg
	if handler, exists := a.msgHandlers[reflect.TypeOf(msg)]; exists {
		if cmd, ok := a.callReflectedFunc(handler, wrappedMsg); ok {
			cmds = append(cmds, cmd)
		}
	}

	// 2. Specific Key Handlers
	slog.Debug("KeyMsg received", "key", msg.String())
	if handler, exists := a.keyHandlers[msg.String()]; exists {
		cmds = append(cmds, unwrapCmd(handler(msg)))
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...), true
	}

	slog.Debug("No controller key handler found", "key", msg.String())
	return nil, false
}

func (a *Application[M]) handleControllerMsg(msg Msg) (tea.Cmd, bool) {
	msgType := reflect.TypeOf(msg.Inner)
	if handler, exists := a.msgHandlers[msgType]; exists {
		slog.Debug("Calling controller message handler", "msg_type", msgType.String())
		if cmd, ok := a.callReflectedFunc(handler, msg); ok {
			return cmd, true
		}
		slog.Debug("Controller message handler returned nil, continuing...")
	}
	return nil, false
}

func (a *Application[M]) handleGlobalKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	for _, handler := range a.globalHandlers {
		if cmd := handler.Handle(msg); cmd != nil {
			slog.Debug("Global handler handled key", "handler", reflect.TypeOf(handler), "key", msg.String())
			return a, cmd
		}
	}
	return a, nil
}

func (a *Application[M]) callReflectedFunc(handler reflect.Value, msg Msg) (tea.Cmd, bool) {
	results := handler.Call([]reflect.Value{reflect.ValueOf(msg.Inner)})
	if len(results) > 0 && !results[0].IsNil() {
		cmd := results[0].Interface().(Cmd)
		unwrapped := unwrapCmd(cmd)
		return unwrapped, true
	}
	return nil, false
}

// gets the appropriate handler based on the paramter
func (a *Application[M]) scanMessageHandlers() {
	slog.Debug("Scanning message handlers")
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

		if strings.HasPrefix(methodName, "On") {
			if methodType.NumIn() == 1 && methodType.NumOut() == 1 {
				msgType := methodType.In(0)
				a.msgHandlers[msgType] = method
			}
		}
	}
	ctlr.Init(a.keyHandlers)
}
