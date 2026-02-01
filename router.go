package mvct

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Route represents a single route
type Route[M any] struct {
	Path       string
	Controller Controller
}

// NavigateMsg signals a route change
type NavigateMsg struct {
	Route string
}

// Router manages routing between controllers
type Router struct {
	routes        map[string]Controller
	currentRoute  string
	previousRoute string
	defaultRoute  string
	middleware    []Middleware
}

func NewRouter() *Router {
	return &Router{
		routes:     make(map[string]Controller),
		middleware: []Middleware{},
	}
}

// Register adds a route
func (r *Router) Register(path string, controller Controller) {
	r.routes[path] = controller
}

// SetDefault sets the default route
func (r *Router) SetDefault(path string) {
	r.defaultRoute = path
	r.currentRoute = path
}

// Use adds middleware
func (r *Router) Use(m Middleware) {
	r.middleware = append(r.middleware, m)
}

// Current returns the current controller
func (r *Router) Current() Controller {
	controller := r.routes[r.currentRoute]
	if controller == nil && r.defaultRoute != "" {
		controller = r.routes[r.defaultRoute]
	}
	return controller
}

// Navigate changes the current route
func (r *Router) Navigate(handlers KeyHandlers, path string) (tea.Cmd, error) {
	if _, ok := r.routes[path]; !ok {
		return nil, fmt.Errorf("route not found: %s", path)
	}

	oldRoute := r.currentRoute
	r.currentRoute = path

	// Run middleware
	ctx := &Context{
		From: oldRoute,
		To:   path,
	}

	for _, m := range r.middleware {
		if !m.Handle(ctx) {
			// Middleware blocked navigation
			r.currentRoute = oldRoute
			return nil, fmt.Errorf("navigation blocked by middleware")
		}
	}

	// Initialize new controller
	return unwrapCmd(r.Current().Init(handlers)), nil
}

// CurrentRoute returns the current route path
func (r *Router) CurrentRoute() string {
	return r.currentRoute
}
