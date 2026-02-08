package mvct

import (
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	handler1 := func(msg KeyMsg) Cmd { return nil }
	handler2 := func(msg KeyMsg) Cmd { return nil }

	// Test building key bindings for single key
	kb1 := Key("a")
	if len(kb1.keys) != 1 || kb1.keys[0] != "a" {
		t.Errorf("Key() failed, expected ['a'], got %v", kb1.keys)
	}

	// Test building key bindings for multiple keys
	kb2 := Key("b", "c")
	if len(kb2.keys) != 2 {
		t.Errorf("Key() for multiple keys failed, expected 2 keys, got %d", len(kb2.keys))
	}

	// Test To handler association
	kb1 = kb1.To(handler1)
	// We can't easily compare function pointers directly in Go tests without reflect,
	// but we can check if it's set.
	if kb1.handler == nil {
		t.Error("To() did not set handler")
	}

	kb2 = kb2.To(handler2)

	// Test Keys map generation
	keyMap := Keys(kb1, kb2)

	if len(keyMap) != 3 { // a, b, c
		t.Errorf("Keys() generated map of wrong size, expected 3, got %d", len(keyMap))
	}

	// Verify mappings
	// reflect.ValueOf(fn).Pointer() can be used to compare functions

	hA := keyMap["a"]
	if reflect.ValueOf(hA).Pointer() != reflect.ValueOf(handler1).Pointer() {
		t.Error("Wrapper for 'a' does not match handler1")
	}

	hB := keyMap["b"]
	if reflect.ValueOf(hB).Pointer() != reflect.ValueOf(handler2).Pointer() {
		t.Error("Wrapper for 'b' does not match handler2")
	}

	hC := keyMap["c"]
	if reflect.ValueOf(hC).Pointer() != reflect.ValueOf(handler2).Pointer() {
		t.Error("Wrapper for 'c' does not match handler2")
	}
}
