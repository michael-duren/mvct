package mvct

// Context provides information about the navigation
type Context struct {
	From string
	To   string
	Data map[string]any
}

// Middleware can intercept route changes
type Middleware interface {
	Handle(ctx *Context) bool // return false to block navigation
}

// MiddlewareFunc is a function adapter for Middleware
type MiddlewareFunc func(ctx *Context) bool

func (f MiddlewareFunc) Handle(ctx *Context) bool {
	return f(ctx)
}

// LoggingMiddleware logs all route changes
func LoggingMiddleware() Middleware {
	return MiddlewareFunc(func(ctx *Context) bool {
		// Could use slog here
		return true
	})
}

// AuthMiddleware checks if routes require authentication
func AuthMiddleware(protectedRoutes []string, isAuthenticated func() bool) Middleware {
	return MiddlewareFunc(func(ctx *Context) bool {
		for _, route := range protectedRoutes {
			if ctx.To == route && !isAuthenticated() {
				return false
			}
		}
		return true
	})
}
