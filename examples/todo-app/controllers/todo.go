package controllers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/bubble-tea-mvc/mvc"
)

type TodoModel struct {
	Todos  []string
	Cursor int
}

type TodosController struct {
	mvc.BaseController
	model TodoModel
}

func NewTodosController() *TodosController {
	return &TodosController{
		model: TodoModel{
			Todos:  []string{},
			Cursor: 0,
		},
	}
}

func (tc *TodosController) Init() tea.Cmd {
	return nil
}

func (tc *TodosController) Update(msg tea.Msg) tea.Cmd {
	model := tc.App().Model().(*AppModel)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if tc.model.Cursor > 0 {
				tc.model.Cursor--
			}
		case "down", "j":
			if tc.model.Cursor < len(model.Todos)-1 {
				tc.model.Cursor++
			}
		case "a":
			model.Todos = append(model.Todos, fmt.Sprintf("Todo #%d", len(model.Todos)+1))
		case "esc":
			return tc.Navigate("/home")
		}
	}
	return nil
}

func (tc *TodosController) View() string {
	model := tc.App().Model().(*AppModel)

	view := "Todos:\n\n"
	for i, todo := range model.Todos {
		cursor := " "
		if i == tc.model.Cursor {
			cursor = ">"
		}
		view += fmt.Sprintf("%s %s\n", cursor, todo)
	}

	view += "\na - Add todo | esc - Back | q - Quit"
	return view
}

func (tc *TodosController) GetModel() TodoModel {
	return tc.model
}
