package mvct

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// Convert tea.KeyMsg to mvct.KeyMsg
func wrapMsg(msg tea.Msg) Msg {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return Msg{
			Inner: KeyMsg{
				Type:  keyMsg.Type,
				Runes: keyMsg.Runes,
				Alt:   keyMsg.Alt,
			},
			Context: context.Background(),
		}
	}

	// Wrap other messages as-is
	return Msg{
		Inner:   msg,
		Context: context.Background(),
	}
}

func wrapCmd(ctx context.Context, cmd tea.Cmd) Cmd {
	if cmd == nil {
		return nil
	}
	return func() Msg {
		return Msg{
			Inner:   cmd(),
			Context: ctx,
		}
	}
}

func unwrapCmd(cmd Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		msg := cmd()
		return msg.Inner
	}
}

func Quit() Cmd {
	return wrapCmd(context.Background(), tea.Quit)
}
