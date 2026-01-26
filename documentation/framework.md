# MVCT Framework - Internal Architecture

## Overview

MVCT is a Model-View-Controller-Template framework built on top of Bubble Tea that provides routing, key handling, and message dispatching without requiring users to write boilerplate Update methods.

## Core Components

### Application[M]

The main orchestrator that manages routing, global state, and message flow.

**Fields:**

- `router *Router` - Routes between controllers
- `model M` - Global application model
- `globalHandlers []GlobalHandler` - Handlers that run before routing (e.g., quit on Ctrl+C)
- `keyHandlers map[string]KeyMsgHandler` - Current controller's registered key handlers
- `msgHandlers map[reflect.Type]reflect.Value` - Current controller's On<MsgType> methods (discovered via reflection)
- `width, height int` - Terminal dimensions

**Responsibilities:**

1. Implements `tea.Model` interface (Init, Update, View)
2. Wraps `tea.Msg` → `mvct.Msg` (adds Context)
3. Unwraps `mvct.Cmd` → `tea.Cmd` (removes Context)
4. Routes messages to current controller's handlers
5. Scans controllers for On<MsgType> methods when routes change
6. Handles global key handlers before routing

### Router

Manages navigation between controllers.

**Fields:**

- `routes map[string]Controller` - Path → Controller mapping
- `currentRoute string` - Active route
- `defaultRoute string` - Fallback route
- `middleware []Middleware` - Pre-navigation hooks

**Responsibilities:**

1. Register controllers at paths
2. Navigate between routes
3. Run middleware on navigation
4. Return current controller

**Navigation Flow:**

1. Receive NavigateMsg
2. Check if route exists
3. Run middleware (can block navigation)
4. Update currentRoute
5. Call new controller's Init()
6. Application scans new controller's handlers

### Controller Interface

Minimal interface users implement.

```go
type Controller interface {
    Init() Cmd
    View() string
    GetModel() any
    RegisterKeyHandlers(handlers map[string]KeyMsgHandler)
}
```

**User Implementation:**

```go
type HomeController struct {
    model HomeModel
}

func (c *HomeController) Init() mvct.Cmd {
    return nil
}

func (c *HomeController) View() string {
    return renderHomeView(c.model)
}

func (c *HomeController) GetModel() any {
    return c.model
}

func (c *HomeController) RegisterKeyHandlers(handlers map[string]mvct.KeyMsgHandler) {
    handlers["q"] = c.onQuit
    handlers["enter"] = c.onEnter
}

// Auto-discovered by reflection
func (c *HomeController) OnWindowSizeMsg(msg mvct.WindowSizeMsg) mvct.Cmd {
    c.model.width = msg.Width
    return nil
}
```

## Message Flow

### 1. Bubble Tea → Application

```
tea.Msg (from Bubble Tea)
  ↓
Application.Update(msg tea.Msg)
```

### 2. Global Handler Check

```
if tea.KeyMsg → run globalHandlers
  if handler returns cmd → return immediately
```

### 3. Message Wrapping

```
tea.Msg → wrapMsg() → mvct.Msg
  - tea.KeyMsg → mvct.KeyMsg
  - Other messages → wrapped as-is
  - Context added
```

### 4. Navigation Handling

```
if NavigateMsg:
  router.Navigate(route)
  scanCurrentControllerHandlers()
  return Init() cmd
```

### 5. Key Handler Dispatch

```
if KeyMsg:
  check keyHandlers map
  if found → call handler(keyMsg)
  return cmd
```

### 6. Message Handler Dispatch (Reflection)

```
if other message type:
  check msgHandlers map (by reflect.Type)
  if found → call method via reflection
  return cmd
```

### 7. Command Unwrapping

```
mvct.Cmd → unwrapCmd() → tea.Cmd
  extracts Inner message
  returns to Bubble Tea
```

## Handler Registration

### Key Handlers (Manual)

User explicitly maps keys to handlers in RegisterKeyHandlers:

```go
func (c *HomeController) RegisterKeyHandlers(handlers map[string]mvct.KeyMsgHandler) {
    handlers["q"] = c.onQuit
    handlers["j"] = c.moveDown
}
```

Application stores these in `keyHandlers` map.

### Message Handlers (Auto-discovered)

Framework scans controller for `On<MsgType>` methods via reflection:

```go
func (c *HomeController) OnWindowSizeMsg(msg mvct.WindowSizeMsg) mvct.Cmd {
    // Auto-registered by framework
}
```

**Scan Process:**

1. On navigation, `scanCurrentControllerHandlers()` runs
2. Reflects on current controller's methods
3. Finds methods starting with "On" (excluding "OnKey")
4. Checks signature: `func(SomeMsg) Cmd`
5. Maps `reflect.TypeOf(SomeMsg)` → method

Application stores these in `msgHandlers` map.

## Helper Functions

### Keys() - Key Binding Builder

```go
// User code
c.RegisterKeyHandlers(mvct.Keys(
    mvct.Key("q", "ctrl+c").To(c.onQuit),
    mvct.Key("j", "down").To(c.moveDown),
))

// Expands to map:
// {
//   "q": c.onQuit,
//   "ctrl+c": c.onQuit,
//   "j": c.moveDown,
//   "down": c.moveDown,
// }
```

## Nested Routing (Future)

Controllers can have their own routers for complex UIs:

```go
type DashboardController struct {
    model DashboardModel
    router *mvct.Router
}

func NewDashboardController() *DashboardController {
    c := &DashboardController{}
    c.router = mvct.NewRouter()
    c.router.Register("/overview", NewOverviewController())
    c.router.Register("/settings", NewSettingsController())
    return c
}
```

**Nested Message Flow:**

1. Message arrives at DashboardController
2. If NavigateMsg with path starting with "/", route to nested controller
3. Otherwise, dispatch to nested router's current controller
4. Parent view wraps child view

## Type Definitions

### Msg

```go
type Msg struct {
    Inner   any              // Actual message (KeyMsg, WindowSizeMsg, custom)
    Context context.Context  // For future: cancellation, tracing, etc.
}
```

### Cmd

```go
type Cmd func() Msg
```

### KeyMsg

```go
type KeyMsg struct {
    Type  tea.KeyType
    Runes []rune
    Alt   bool
}

func (k KeyMsg) String() string // Returns "enter", "ctrl+c", etc.
```

### KeyMsgHandler

```go
type KeyMsgHandler func(msg KeyMsg) Cmd
```

## Initialization Flow

```
1. User creates Application:
   app := mvct.NewApplication(config, globalModel)

2. User registers controllers:
   app.RegisterController("/home", &HomeController{})
   app.RegisterController("/settings", &SettingsController{})

3. User registers global handlers:
   app.UseGlobalHandler(mvct.QuitHandler("ctrl+c", "q"))

4. User runs application:
   app.Run()

5. Bubble Tea calls Init():
   router.Current().Init()
   scanCurrentControllerHandlers()

6. Message loop begins
```

## Navigation Flow

```
1. Controller returns NavigateMsg:
   return func() Msg { return Msg{Inner: NavigateMsg{Route: "/settings"}} }

2. Application receives it in Update:
   if navMsg, ok := wrappedMsg.Inner.(NavigateMsg)

3. Router.Navigate() called:
   - Validate route exists
   - Run middleware
   - Update currentRoute
   - Return new controller's Init()

4. Application scans new controller:
   scanCurrentControllerHandlers()
   - Clear old keyHandlers and msgHandlers
   - Call controller.RegisterKeyHandlers()
   - Reflect for On<MsgType> methods

5. New controller is active
```

## Design Principles

1. **Zero boilerplate** - Users don't write Update methods
2. **Convention over configuration** - `On<MsgType>` methods auto-discovered
3. **Explicit key bindings** - Keys must be registered in RegisterKeyHandlers
4. **Type safety** - Generics for models, reflection only where needed
5. **Composability** - Controllers can have nested routers
6. **Simplicity** - Minimal interface, framework handles complexity

## What Users Don't See

- Message wrapping/unwrapping
- Reflection for method scanning
- Internal routing logic
- Command conversion between tea.Cmd and mvct.Cmd
- Handler map management

## What Users Do

1. Implement 4 interface methods
2. Write handler methods (key handlers + On<MsgType> handlers)
3. Register key bindings
4. Return NavigateMsg to change routes
5. Return other Cmds for async operations
