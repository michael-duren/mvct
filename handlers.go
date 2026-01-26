package mvct

import tea "github.com/charmbracelet/bubbletea"

// GlobalHandler handles messages before routing
type GlobalHandler interface {
	Handle(msg tea.KeyMsg) tea.Cmd
}

// GlobalHandlerFunc is a function adapter
type GlobalHandlerFunc func(msg tea.KeyMsg) tea.Cmd

func (f GlobalHandlerFunc) Handle(msg tea.KeyMsg) tea.Cmd {
	return f(msg)
}

// QuitHandler provides global quit functionality
func QuitHandler(keys ...string) GlobalHandler {
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	return GlobalHandlerFunc(func(msg tea.KeyMsg) tea.Cmd {
		if keyMap[msg.String()] {
			return tea.Quit
		}
		return nil
	})
}

type KeyMsg tea.KeyMsg
