package controllers

import (
	"github.com/michael-duren/mvct"
)

type ExitController struct{}

func NewExitController() *ExitController {
	return &ExitController{}
}

func (c *ExitController) Init(handlers mvct.KeyHandlers) mvct.Cmd {
	return mvct.Quit()
}

func (c *ExitController) View() string {
	return "Goodbye!\n"
}

func (c *ExitController) GetModel() any {
	return nil
}
