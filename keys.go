package mvct

// KeyMsgHandler associates a list of keys with a handler function
type KeyBinding struct {
	keys    []string
	handler KeyMsgHandler
}

// Key creates a new key binding with the specified keys
func Key(keys ...string) KeyBinding {
	return KeyBinding{
		keys: keys,
	}
}

// To assigns a handler to the key binding
func (kb KeyBinding) To(handler KeyMsgHandler) KeyBinding {
	kb.handler = handler
	return kb
}

// Keys converts a list of KeyBindings into a map of key strings to handlers
// this is registered with a controller
func Keys(bindings ...KeyBinding) map[string]KeyMsgHandler {
	keyMap := make(map[string]KeyMsgHandler)
	for _, binding := range bindings {
		for _, k := range binding.keys {
			keyMap[k] = binding.handler
		}
	}
	return keyMap
}
