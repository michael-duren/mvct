package controllers

import (
	todoapp "example-todo"
	"fmt"

	"github.com/michael-duren/mvct"
)

type TodoModel struct {
	Todos  []string
	Cursor int
}

type TodosController struct {
	mvct.BaseController[todoapp.AppModel]
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

func (tc *TodosController) Init() mvct.Cmd {
	return nil
}

func (tc *TodosController) Update(msg mvct.Msg) {

}

func (tc *TodosController) View() string {
	view := "Todos:\n\n"
	for i, todo := range tc.model.Todos {
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
