package todoapp

import "github.com/michael-duren/bubble-tea-mvc/mvc"

type AppModel struct {
}

func main() {
	appModel := &AppModel{}
	app := mvc.NewApplication[*AppModel](mvc.Config{
		DefaultRoute: "home",
	}, appModel)
	// mvc.RegisterController(app, "home", NewHomeController())
	// app.Use()
}

// // examples/todo-app/main.go
// package main
//
// import (
// 	"fmt"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/yourname/bubbletea-mvc/framework"
// )
//
// // Your app model
// type AppModel struct {
// 	Todos []string
// 	User  *User
// }
//
// type User struct {
// 	Name          string
// 	Authenticated bool
// }
//
// func main() {
// 	// Create your app model
// 	model := &AppModel{
// 		Todos: []string{},
// 		User:  &User{Name: "Guest", Authenticated: false},
// 	}
//
// 	// Create application
// 	app := framework.NewApplication(framework.Config{
// 		DefaultRoute: "/home",
// 		Model:        model,
// 	})
//
// 	// Register global handlers
// 	app.UseGlobalHandler(framework.QuitHandler("ctrl+c", "q"))
//
// 	// Register middleware
// 	app.Use(framework.LoggingMiddleware())
// 	app.Use(framework.AuthMiddleware(
// 		[]string{"/settings"},
// 		func() bool {
// 			m := app.Model().(*AppModel)
// 			return m.User.Authenticated
// 		},
// 	))
//
// 	// Register controllers
// 	app.RegisterController("/home", NewHomeController())
// 	app.RegisterController("/todos", NewTodosController())
// 	app.RegisterController("/settings", NewSettingsController())
//
// 	// Run
// 	p := tea.NewProgram(app)
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// }
//
// // HomeController
// type HomeController struct {
// 	framework.BaseController
// }
//
// func NewHomeController() *HomeController {
// 	return &HomeController{}
// }
//
// func (hc *HomeController) Init() tea.Cmd {
// 	return nil
// }
//
// func (hc *HomeController) Update(msg tea.Msg) tea.Cmd {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "t":
// 			return hc.Navigate("/todos")
// 		case "s":
// 			return hc.Navigate("/settings")
// 		}
// 	}
// 	return nil
// }
//
// func (hc *HomeController) View() string {
// 	model := hc.App().Model().(*AppModel)
//
// 	return fmt.Sprintf(`
// Welcome, %s!
//
// Commands:
//   t - View todos
//   s - Settings
//   q - Quit
// 	`, model.User.Name)
// }
//
// // TodosController
// type TodosController struct {
// 	framework.BaseController
// 	cursor int
// }
//
// func NewTodosController() *TodosController {
// 	return &TodosController{cursor: 0}
// }
//
// func (tc *TodosController) Init() tea.Cmd {
// 	return nil
// }
//
// func (tc *TodosController) Update(msg tea.Msg) tea.Cmd {
// 	model := tc.App().Model().(*AppModel)
//
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "up", "k":
// 			if tc.cursor > 0 {
// 				tc.cursor--
// 			}
// 		case "down", "j":
// 			if tc.cursor < len(model.Todos)-1 {
// 				tc.cursor++
// 			}
// 		case "a":
// 			model.Todos = append(model.Todos, fmt.Sprintf("Todo #%d", len(model.Todos)+1))
// 		case "esc":
// 			return tc.Navigate("/home")
// 		}
// 	}
// 	return nil
// }
//
// func (tc *TodosController) View() string {
// 	model := tc.App().Model().(*AppModel)
//
// 	view := "Todos:\n\n"
// 	for i, todo := range model.Todos {
// 		cursor := " "
// 		if i == tc.cursor {
// 			cursor = ">"
// 		}
// 		view += fmt.Sprintf("%s %s\n", cursor, todo)
// 	}
//
// 	view += "\na - Add todo | esc - Back | q - Quit"
// 	return view
// }
//
// // SettingsController
// type SettingsController struct {
// 	framework.BaseController
// }
//
// func NewSettingsController() *SettingsController {
// 	return &SettingsController{}
// }
//
// func (sc *SettingsController) Init() tea.Cmd {
// 	return nil
// }
//
// func (sc *SettingsController) Update(msg tea.Msg) tea.Cmd {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "esc":
// 			return sc.Navigate("/home")
// 		}
// 	}
// 	return nil
// }
//
// func (sc *SettingsController) View() string {
// 	return "Settings\n\nesc - Back | q - Quit"
// }
