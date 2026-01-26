package controllers

import (
	"example-todo/components"
	"fmt"
	"strings"

	"github.com/michael-duren/mvct"
)

type TodoModel struct {
	todos  []string
	cursor int
	input  string
	adding bool
}

type TodoController struct {
	model TodoModel
}

func NewTodoController() *TodoController {
	return &TodoController{
		model: TodoModel{
			todos: []string{
				"Buy groceries",
				"Walk the dog",
				"Write code",
			},
		},
	}
}

func (c *TodoController) Init() mvct.Cmd {
	return nil
}

func (c *TodoController) View() string {
	var b strings.Builder

	if c.model.adding {
		b.WriteString(components.HeaderStyle.Render("Add New Todo"))
		b.WriteString("\n\n")
		b.WriteString(components.InputStyle.Render(c.model.input + "█"))
		b.WriteString("\n\n")
		b.WriteString(components.HelpStyle.Render("enter: save • esc: cancel"))
		return b.String()
	}

	// Todo list
	b.WriteString(components.HeaderStyle.Render(fmt.Sprintf("Tasks (%d)", len(c.model.todos))))
	b.WriteString("\n\n")

	if len(c.model.todos) == 0 {
		b.WriteString(components.UnselectedItemStyle.Render("  No todos yet. Press 'a' to add one!"))
	} else {
		for i, todo := range c.model.todos {
			if c.model.cursor == i {
				b.WriteString(components.SelectedItemStyle.Render("▶ " + todo))
			} else {
				b.WriteString(components.UnselectedItemStyle.Render("  " + todo))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(components.HelpStyle.Render("↑/k: up • ↓/j: down • a: add • d: delete • q: quit"))

	return b.String()
}

func (c *TodoController) GetModel() any {
	return c.model
}

func (c *TodoController) RegisterKeyHandlers(handlers map[string]mvct.KeyMsgHandler) {
	for k, v := range mvct.Keys(
		mvct.Key("q").To(c.onQuit),
		mvct.Key("j", "down").To(c.onDown),
		mvct.Key("k", "up").To(c.onUp),
		mvct.Key("a").To(c.onAdd),
		mvct.Key("d").To(c.onDelete),
		mvct.Key("enter").To(c.onEnter),
		mvct.Key("esc").To(c.onEscape),
	) {
		handlers[k] = v
	}
}

func (c *TodoController) onQuit(msg mvct.KeyMsg) mvct.Cmd {
	return mvct.Navigate(R.Exit)
}

func (c *TodoController) onDown(msg mvct.KeyMsg) mvct.Cmd {
	if !c.model.adding && c.model.cursor < len(c.model.todos)-1 {
		c.model.cursor++
	}
	return nil
}

func (c *TodoController) onUp(msg mvct.KeyMsg) mvct.Cmd {
	if !c.model.adding && c.model.cursor > 0 {
		c.model.cursor--
	}
	return nil
}

func (c *TodoController) onAdd(msg mvct.KeyMsg) mvct.Cmd {
	if !c.model.adding {
		c.model.adding = true
		c.model.input = ""
	}
	return nil
}

func (c *TodoController) onDelete(msg mvct.KeyMsg) mvct.Cmd {
	if !c.model.adding && len(c.model.todos) > 0 {
		c.model.todos = append(
			c.model.todos[:c.model.cursor],
			c.model.todos[c.model.cursor+1:]...,
		)
		if c.model.cursor >= len(c.model.todos) && c.model.cursor > 0 {
			c.model.cursor--
		}
	}
	return nil
}

func (c *TodoController) onEnter(msg mvct.KeyMsg) mvct.Cmd {
	if c.model.adding && c.model.input != "" {
		c.model.todos = append(c.model.todos, c.model.input)
		c.model.adding = false
		c.model.input = ""
	}
	return nil
}

func (c *TodoController) onEscape(msg mvct.KeyMsg) mvct.Cmd {
	if c.model.adding {
		c.model.adding = false
		c.model.input = ""
	}
	return nil
}

// Handle text input when adding
func (c *TodoController) OnKeyMsg(msg mvct.KeyMsg) mvct.Cmd {
	if c.model.adding {
		// Handle typing
		if len(msg.Runes) > 0 {
			c.model.input += string(msg.Runes)
		}
		// Handle backspace
		if msg.Type == mvct.KeyBackspace && len(c.model.input) > 0 {
			c.model.input = c.model.input[:len(c.model.input)-1]
		}
	}
	return nil
}
