package mvct

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type TestMsg struct{ content string }

func TestWrapMsg(t *testing.T) {
	// Test wrapping custom message
	orig := TestMsg{content: "test"}
	wrapper := wrapMsg(orig)

	if wrapper.Inner != orig {
		t.Error("wrapMsg did not preserve inner custom message")
	}
	if wrapper.Context == nil {
		t.Error("wrapMsg did not create context")
	}

	// Test wrapping KeyMsg
	km := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a"), Alt: true}
	wrapperKm := wrapMsg(km)

	innerKm, ok := wrapperKm.Inner.(KeyMsg)
	if !ok {
		t.Error("wrapMsg did not convert tea.KeyMsg to mvct.KeyMsg")
	}
	if innerKm.Type != km.Type || innerKm.Alt != km.Alt {
		t.Error("wrapMsg did not preserve KeyMsg fields")
	}
}

func TestWrapUnwrapCmd(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "val")

	// Create a tea.Cmd that returns a message
	origMsg := TestMsg{content: "cmd"}
	teaCmd := func() tea.Msg {
		return origMsg
	}

	// Wrap it
	cmd := wrapCmd(ctx, teaCmd)
	if cmd == nil {
		t.Fatal("wrapCmd returned nil")
	}

	// Execute wrapped cmd to get Msg
	msg := cmd()
	if msg.Inner != origMsg {
		t.Error("Wrapped cmd did not return original inner message")
	}
	if msg.Context.Value("key") != "val" {
		t.Error("Wrapped cmd did not preserve context")
	}

	// Unwrap it back to tea.Cmd
	unwrapped := unwrapCmd(cmd)
	if unwrapped == nil {
		t.Fatal("unwrapCmd returned nil")
	}

	// Execute unwrapped cmd to get tea.Msg
	teaMsg := unwrapped()
	if teaMsg != origMsg {
		t.Error("Unwrapped cmd did not return original message")
	}
}

func TestQuit(t *testing.T) {
	cmd := Quit()
	if cmd == nil {
		t.Fatal("Quit returned nil")
	}

	msg := cmd()
	if _, ok := msg.Inner.(tea.QuitMsg); !ok {
		t.Error("Quit cmd did not return QuitMsg")
	}
}
