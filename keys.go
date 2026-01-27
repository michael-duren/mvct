package mvct

import tea "github.com/charmbracelet/bubbletea"

type KeyMsg = tea.KeyMsg

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

const (
	KeyBackspace = tea.KeyBackspace
	KeyEnter     = tea.KeyEnter
	KeyEsc       = tea.KeyEsc
	KeySpace     = tea.KeySpace
	KeyTab       = tea.KeyTab
	KeyShiftTab  = tea.KeyShiftTab
	KeyUp        = tea.KeyUp
	KeyDown      = tea.KeyDown
	KeyLeft      = tea.KeyLeft
	KeyRight     = tea.KeyRight
	KeyHome      = tea.KeyHome
	KeyEnd       = tea.KeyEnd
	KeyPgUp      = tea.KeyPgUp
	KeyPgDown    = tea.KeyPgDown
	KeyDelete    = tea.KeyDelete
	KeyInsert    = tea.KeyInsert
	KeyCtrlA     = tea.KeyCtrlA
	KeyCtrlB     = tea.KeyCtrlB
	KeyCtrlC     = tea.KeyCtrlC
	KeyCtrlD     = tea.KeyCtrlD
	KeyCtrlE     = tea.KeyCtrlE
	KeyCtrlF     = tea.KeyCtrlF
	KeyCtrlG     = tea.KeyCtrlG
	KeyCtrlH     = tea.KeyCtrlH
	KeyCtrlI     = tea.KeyCtrlI
	KeyCtrlJ     = tea.KeyCtrlJ
	KeyCtrlK     = tea.KeyCtrlK
	KeyCtrlL     = tea.KeyCtrlL
	KeyCtrlM     = tea.KeyCtrlM
	KeyCtrlN     = tea.KeyCtrlN
	KeyCtrlO     = tea.KeyCtrlO
	KeyCtrlP     = tea.KeyCtrlP
	KeyCtrlQ     = tea.KeyCtrlQ
	KeyCtrlR     = tea.KeyCtrlR
	KeyCtrlS     = tea.KeyCtrlS
	KeyCtrlT     = tea.KeyCtrlT
	KeyCtrlU     = tea.KeyCtrlU
	KeyCtrlV     = tea.KeyCtrlV
	KeyCtrlW     = tea.KeyCtrlW
	KeyCtrlX     = tea.KeyCtrlX
	KeyCtrlY     = tea.KeyCtrlY
	KeyCtrlZ     = tea.KeyCtrlZ
	KeyF1        = tea.KeyF1
	KeyF2        = tea.KeyF2
	KeyF3        = tea.KeyF3
	KeyF4        = tea.KeyF4
	KeyF5        = tea.KeyF5
	KeyF6        = tea.KeyF6
	KeyF7        = tea.KeyF7
	KeyF8        = tea.KeyF8
	KeyF9        = tea.KeyF9
	KeyF10       = tea.KeyF10
	KeyF11       = tea.KeyF11
	KeyF12       = tea.KeyF12
	KeyF13       = tea.KeyF13
	KeyF14       = tea.KeyF14
	KeyF15       = tea.KeyF15
	KeyF16       = tea.KeyF16
	KeyF17       = tea.KeyF17
	KeyF18       = tea.KeyF18
	KeyF19       = tea.KeyF19
	KeyF20       = tea.KeyF20
)
