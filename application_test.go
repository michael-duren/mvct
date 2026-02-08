package mvct

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type StringMsg struct {
	Value string
}

type ReflectionController struct {
	MockController
	receivedMsg string
	keyHandled  bool
}

func (c *ReflectionController) OnStringMsg(msg StringMsg) Cmd {
	c.receivedMsg = msg.Value
	return nil
}

func (c *ReflectionController) Init(handlers KeyHandlers) Cmd {
	handlers["x"] = func(msg KeyMsg) Cmd {
		c.keyHandled = true
		return nil
	}
	return nil
}

func TestNewApplication(t *testing.T) {
	config := Config{
		DefaultRoute: "/home",
	}
	model := "test_model"
	app := NewApplication(config, model)

	if app == nil {
		t.Fatal("NewApplication returned nil")
	}
	if app.router.defaultRoute != "/home" {
		t.Errorf("expected default route /home, got %s", app.router.defaultRoute)
	}
	if app.Model() != model {
		t.Errorf("expected model %s, got %v", model, app.Model())
	}
}

func TestApplicationInit(t *testing.T) {
	config := Config{DefaultRoute: "/init"}
	app := NewApplication(config, "model")
	c := &MockController{name: "init_controller"}
	app.RegisterController("/init", c)

	_ = app.Init()
	if !c.initCalled {
		t.Error("Init did not call controller Init")
	}
	// cmd can be nil if controller Init returns nil
}

func TestApplicationUpdate_WindowSize(t *testing.T) {
	app := NewApplication(Config{}, "model")
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	_, _ = app.Update(msg)

	if app.Width() != 100 {
		t.Errorf("expected width 100, got %d", app.Width())
	}
	if app.Height() != 50 {
		t.Errorf("expected height 50, got %d", app.Height())
	}
}

func TestApplicationUpdate_Navigation(t *testing.T) {
	app := NewApplication(Config{DefaultRoute: "/start"}, "model")
	c1 := &MockController{name: "start"}
	c2 := &MockController{name: "next"}
	app.RegisterController("/start", c1)
	app.RegisterController("/next", c2)

	// Trigger initial init to check reflected handlers later
	app.Init()

	navMsg := NavigateMsg{Route: "/next"}
	// Update usually takes tea.Msg, our NavigateMsg is wrapped inside Msg logic inside Update
	// But Update takes tea.Msg. Wait, application.go Update takes tea.Msg.
	// NavigateMsg is internal. How does it get there?
	// It comes from a command. Bubble tea loop sends it.
	// But Update expects tea.Msg. NavigateMsg is NOT tea.Msg (it's interface{}).
	// But in `update.go`:
	// switch msg := msg.(type) ...
	// wrappedMsg := wrapMsg(msg)
	// switch inner := wrappedMsg.Inner.(type) { case NavigateMsg: ... }
	// So passing NavigateMsg directly to Update should work if it's treated as a custom Msg?
	// Wait, NavigateMsg is struct. tea.Msg is any.
	// The implementation in update.go wraps it. wrapMsg handles tea.KeyMsg.
	// If I pass NavigateMsg, wrapMsg returns Msg{Inner: NavigateMsg}.
	// Then switch inner matches NavigateMsg.

	app.Update(navMsg)

	if app.router.CurrentRoute() != "/next" {
		t.Errorf("expected route /next, got %s", app.router.CurrentRoute())
	}
	if !c2.initCalled {
		t.Error("Navigate did not init new controller")
	}
}

func TestApplicationUpdate_ReflectionHandler(t *testing.T) {
	app := NewApplication(Config{DefaultRoute: "/refl"}, "model")
	rc := &ReflectionController{}
	app.RegisterController("/refl", rc)
	app.Init() // Scans handlers

	msg := StringMsg{Value: "hello"}
	app.Update(msg)

	if rc.receivedMsg != "hello" {
		t.Errorf("expected reflected handler to receive 'hello', got '%s'", rc.receivedMsg)
	}
}

func TestApplicationUpdate_KeyHandler(t *testing.T) {
	app := NewApplication(Config{DefaultRoute: "/key"}, "model")
	rc := &ReflectionController{}
	app.RegisterController("/key", rc)
	app.Init() // Registers 'x' handler

	// Test controller key handler
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
	app.Update(keyMsg)

	if !rc.keyHandled {
		t.Error("Controller key handler for 'x' was not called")
	}

	// Test global handler
	globalHandled := false
	app.UseGlobalHandler(GlobalHandlerFunc(func(msg tea.KeyMsg) tea.Cmd {
		if msg.String() == "q" {
			globalHandled = true
			return tea.Quit
		}
		return nil
	}))

	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := app.Update(quitMsg)

	if !globalHandled {
		t.Error("Global handler was not called")
	}
	if cmd == nil {
		t.Error("Global handler should return Quit cmd")
	}
}
