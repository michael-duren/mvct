package templates

const GoModTemplate = `module {{.Name}}

go 1.25.1

require (
	github.com/michael-duren/mvct v0.0.0
	github.com/charmbracelet/log v0.4.0
)

replace github.com/michael-duren/mvct => ../
`

const MainGoTemplate = `package main

import (
	"log/slog"
{{if .WithDB}}
    "database/sql"
    "log"
	_ "github.com/mattn/go-sqlite3"
{{end}}

	"github.com/michael-duren/mvct"
	"github.com/charmbracelet/log"
	"{{.Name}}/components"
	"{{.Name}}/controllers"
)

type AppModel struct {
{{if .WithDB}}
    DB *sql.DB
{{end}}
}

func main() {
{{if .WithDB}}
    db, err := sql.Open("sqlite3", "app.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
{{end}}

	R := controllers.R
	app := mvct.NewApplication(mvct.Config{
		DefaultRoute: R.Home,
	}, AppModel{
{{if .WithDB}}
        DB: db,
{{end}}
    })

	app.UseLogger(mvct.LoggerConfig{
		Path:     "app.log",
		LogLevel: mvct.LogLevel(slog.LevelDebug),
	})

	app.SetLayout(func(content string, width, height int) string {
		return components.Layout(content, width, height)
	})

	app.UseGlobalHandler(mvct.QuitHandler("ctrl+c"))
	app.UseGlobalHandler(mvct.QuitHandler("q"))

	app.RegisterController(R.Home, controllers.NewHomeController())

	app.Run()
}
`

const LayoutTemplate = `package components

import "github.com/charmbracelet/lipgloss"

func Layout(content string, width, height int) string {
    style := lipgloss.NewStyle().
        Width(width).
        Height(height).
        Align(lipgloss.Center, lipgloss.Center).
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62"))

    return style.Render(content)
}
`

const HomeControllerTemplate = `package controllers

import (
    "fmt"
	"github.com/michael-duren/mvct"
	container "github.com/charmbracelet/bubbletea"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (c *HomeController) Init(m mvct.Model) (mvct.Model, container.Cmd) {
	return m, nil
}

func (c *HomeController) Update(m mvct.Model, msg container.Msg) (mvct.Model, container.Cmd) {
    if msg, ok := msg.(container.KeyMsg); ok {
        if msg.String() == "enter" {
            // Do something
        }
    }
	return m, nil
}

func (c *HomeController) View(m mvct.Model) string {
	return fmt.Sprintf("Welcome to %s!\nPress q to quit.", "{{.Name}}")
}
`

const RoutingTemplate = `package controllers

import "github.com/michael-duren/mvct"

var R = struct {
    Home mvct.Route
}{
    Home: mvct.Route("home"),
}
`

const SqlcYamlTemplate = `version: "2"
sql:
  - schema: "db/schema.sql"
    queries: "db/queries.sql"
    engine: "sqlite"
    gen:
      go:
        package: "db"
        out: "internal/db"
`

const SchemaSqlTemplate = `CREATE TABLE todos (
  id   INTEGER PRIMARY KEY,
  text TEXT    NOT NULL,
  done BOOLEAN NOT NULL DEFAULT 0
);
`

const QueriesSqlTemplate = `-- name: GetTodos :many
SELECT * FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (text, done) VALUES (?, ?) RETURNING *;

-- name: UpdateTodo :exec
UPDATE todos SET done = ? WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;
`

const MakefileTemplate = `sqlc:
	sqlc generate
`
