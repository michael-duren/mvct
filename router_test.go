package mvct

import (
	"testing"
)

// MockController is a simple controller for testing
type MockController struct {
	initCalled bool
	viewCalled bool
	name       string
}

func (c *MockController) Init(handlers KeyHandlers) Cmd {
	c.initCalled = true
	return nil
}

func (c *MockController) View() string {
	c.viewCalled = true
	return c.name
}

// MockMiddleware is a simple middleware for testing
type MockMiddleware struct {
	handleCalled bool
	shouldBlock  bool
}

func (m *MockMiddleware) Handle(ctx *Context) bool {
	m.handleCalled = true
	return !m.shouldBlock
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Fatal("NewRouter returned nil")
	}
	if router.routes == nil {
		t.Error("Router routes map is nil")
	}
	if router.middleware == nil {
		t.Error("Router middleware slice is nil")
	}
}

func TestRegisterAndCurrent(t *testing.T) {
	router := NewRouter()
	c1 := &MockController{name: "c1"}
	router.Register("/c1", c1)

	// Test Register
	if _, ok := router.routes["/c1"]; !ok {
		t.Error("Register failed to add route")
	}

	// Test Current with no current route and no default
	if router.Current() != nil {
		t.Error("Current should return nil when no route is set")
	}

	// Test SetDefault
	router.SetDefault("/c1")
	if router.defaultRoute != "/c1" {
		t.Errorf("expected default route to be /c1, got %s", router.defaultRoute)
	}
	if router.Current() != c1 {
		t.Error("Current should return default route controller")
	}
}
func TestNavigate(t *testing.T) {
	router := NewRouter()
	c1 := &MockController{name: "c1"}
	c2 := &MockController{name: "c2"}

	router.Register("/c1", c1)
	router.Register("/c2", c2)
	router.SetDefault("/c1")

	// Test Navigate to existing route
	cmd, err := router.Navigate(nil, "/c2")
	if err != nil {
		t.Errorf("Navigate failed: %v", err)
	}
	if router.CurrentRoute() != "/c2" {
		t.Errorf("expected current route to be /c2, got %s", router.CurrentRoute())
	}
	if !c2.initCalled {
		t.Error("Navigate did not call Init on new controller")
	}
	// cmd is a function wrapping another function, difficult to test equality directly without execution,
	// but we can check it's not nil if Init returned something (Mock returns nil though)
	_ = cmd

	// Test Navigate to non-existing route
	_, err = router.Navigate(nil, "/missing")
	if err == nil {
		t.Error("Navigate should fail for non-existing route")
	}
	if router.CurrentRoute() != "/c2" {
		t.Errorf("current route should remain /c2 after failed navigation, got %s", router.CurrentRoute())
	}
}

func TestMiddleware(t *testing.T) {
	router := NewRouter()
	c1 := &MockController{name: "c1"}
	c2 := &MockController{name: "c2"}
	router.Register("/c1", c1)
	router.Register("/c2", c2)
	router.SetDefault("/c1")

	// Test allowing middleware
	mwAllow := &MockMiddleware{shouldBlock: false}
	router.Use(mwAllow)

	_, err := router.Navigate(nil, "/c2")
	if err != nil {
		t.Errorf("Navigate should succeed with allowing middleware: %v", err)
	}
	if !mwAllow.handleCalled {
		t.Error("Middleware Handle was not called")
	}
	if router.CurrentRoute() != "/c2" {
		t.Errorf("expected route /c2, got %s", router.CurrentRoute())
	}

	// Test blocking middleware
	mwBlock := &MockMiddleware{shouldBlock: true}
	router.Use(mwBlock)

	_, err = router.Navigate(nil, "/c1")
	if err == nil {
		t.Error("Navigate should fail with blocking middleware")
	}
	if !mwBlock.handleCalled {
		t.Error("Blocking Middleware Handle was not called")
	}
	if router.CurrentRoute() != "/c2" {
		t.Errorf("route should remain /c2 after blocked navigation, got %s", router.CurrentRoute())
	}
}
