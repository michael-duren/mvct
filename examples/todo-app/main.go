package main

import (
	"example-todo/components"
	"example-todo/controllers"

	"github.com/michael-duren/mvct"
)

type AppModel struct {
}

func main() {
	R := controllers.R
	app := mvct.NewApplication(mvct.Config{
		DefaultRoute: R.Home,
	}, AppModel{})

	app.SetLayout(func(content string, width, height int) string {
		layout := components.NewLayout("📝 TODO APP")
		layout.Content = content
		layout.Footer = "Press ? for help"
		layout.Width = width
		layout.Height = height
		return layout.Render()
	})

	app.UseGlobalHandler(mvct.QuitHandler("ctrl+c"))
	app.UseGlobalHandler(mvct.QuitHandler("q"))

	app.UseLogger(mvct.LoggerConfig{
		Path: "",
	})

	app.RegisterController(R.Home, controllers.NewTodoController())
	app.RegisterController(R.Exit, controllers.NewExitController())

	app.Run()
}
