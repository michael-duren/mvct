package mvct

import (
	"context"
	"reflect"
	"strings"
)

// BaseController provides common functionality
// and global access to the Application instance
type BaseController[M any] struct {
	app *Application[M]
	// key handlers are registered via the RegisterKeyHandlers method and are mapped by key string, use convience methods to create
	// these maps
	keyHandlers map[string]KeyMsgHandler
	// message handlers are registered by type - auto-discovered via reflection
	msgHandlers map[reflect.Type]reflect.Value
	// optional router for nested views
	router *Router
}

func (bc *BaseController[M]) SetApp(app *Application[M]) {
	bc.app = app
	if bc.keyHandlers == nil {
		bc.keyHandlers = make(map[string]KeyMsgHandler)
	}
	if bc.msgHandlers == nil {
		bc.msgHandlers = make(map[reflect.Type]reflect.Value)
	}
	bc.scanMessageHandlers()
}

func (bc *BaseController[M]) App() *Application[M] {
	return bc.app
}

// Navigate changes the current route
func (bc *BaseController[M]) Navigate(route string) Cmd {
	return func() Msg {
		return Msg{
			Inner:   NavigateMsg{Route: route},
			Context: context.Background(),
		}
	}
}

// NavigateMsg signals a route change
type NavigateMsg struct {
	Route string
}

func (bc *BaseController[M]) internalUpdate(msg Msg) Cmd {
	if keyMsg, ok := msg.Inner.(KeyMsg); ok {
		if handler, exists := bc.keyHandlers[keyMsg.Type.String()]; exists {
			return handler(keyMsg)
		}
	}

	msgType := reflect.TypeOf(msg.Inner)
	if handler, exists := bc.msgHandlers[msgType]; exists {
		results := handler.Call([]reflect.Value{reflect.ValueOf(msg.Inner)})
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().(Cmd)
		}
	}
	return nil
}

// gets the appropriate handler based on the paramter
func (bc *BaseController[M]) scanMessageHandlers() {
	val := reflect.ValueOf(bc)
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
				bc.msgHandlers[msgType] = method
			}
		}
	}
}
