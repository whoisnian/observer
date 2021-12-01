//go:generate stringer -type=Key -trimprefix=K_ -output=key_name.go
package driver

type Key = byte
type Keycodes [3]Key

const (
	K_ESC Key = iota + 1

	K_F1
	K_F2
	K_F3
	K_F4
	K_F5
	K_F6
	K_F7
	K_F8
	K_F9
	K_F10
	K_F11
	K_F12

	K_0
	K_1
	K_2
	K_3
	K_4
	K_5
	K_6
	K_7
	K_8
	K_9

	K_A
	K_B
	K_C
	K_D
	K_E
	K_F
	K_G
	K_H
	K_I
	K_J
	K_K
	K_L
	K_M
	K_N
	K_O
	K_P
	K_Q
	K_R
	K_S
	K_T
	K_U
	K_V
	K_W
	K_X
	K_Y
	K_Z

	K_GRAVE      // `
	K_MINUS      // -
	K_EQUAL      // =
	K_BACKSPACE  // Backspace
	K_TAB        // Tab
	K_LEFTBRACE  // [
	K_RIGHTBRACE // ]
	K_BACKSLASH  // \
	K_SEMICOLON  // ;
	K_APOSTROPHE // '
	K_ENTER      // Enter
	K_SHIFT      // Shift
	K_COMMA      // ,
	K_DOT        // .
	K_SLASH      // /
	K_CTRL       // Ctrl
	K_ALT        // Alt
	K_SPACE      // Space
	K_INSERT     // Insert
	K_DELETE     // Delete
	K_HOME       // Home
	K_END        // End
	K_PAGEUP     // Page Up
	K_PAGEDOWN   // Page Down
	K_UP         // Up
	K_DOWN       // Down
	K_LEFT       // Left
	K_RIGHT      // Right
)

func (ks Keycodes) String() string {
	if ks[1] == 0 && ks[2] == 0 {
		return KeyName(ks[0])
	} else if ks[2] == 0 {
		return KeyName(ks[0]) + " + " + KeyName(ks[1])
	} else {
		return KeyName(ks[0]) + " + " + KeyName(ks[1]) + " + " + KeyName(ks[2])
	}
}

func VT100Decode(b []byte) Keycodes {
	if len(b) == 1 {
		return map[byte]Keycodes{
			0x00: {K_CTRL, K_2},
			0x01: {K_CTRL, K_A},
			0x02: {K_CTRL, K_B},
			0x03: {K_CTRL, K_C},
			0x04: {K_CTRL, K_D},
			0x05: {K_CTRL, K_E},
			0x06: {K_CTRL, K_F},
			0x07: {K_CTRL, K_G},
			0x08: {K_CTRL, K_H},
			0x09: {K_CTRL, K_I}, // K_TAB
			0x0a: {K_CTRL, K_J},
			0x0b: {K_CTRL, K_K},
			0x0c: {K_CTRL, K_L},
			0x0d: {K_CTRL, K_M}, // K_ENTER
			0x0e: {K_CTRL, K_N},
			0x0f: {K_CTRL, K_O},
			0x10: {K_CTRL, K_P},
			0x11: {K_CTRL, K_Q},
			0x12: {K_CTRL, K_R},
			0x13: {K_CTRL, K_S},
			0x14: {K_CTRL, K_T},
			0x15: {K_CTRL, K_U},
			0x16: {K_CTRL, K_V},
			0x17: {K_CTRL, K_W},
			0x18: {K_CTRL, K_X},
			0x19: {K_CTRL, K_Y},
			0x1a: {K_CTRL, K_Z},
			0x1b: {K_CTRL, K_3}, // K_ESC
			0x1c: {K_CTRL, K_4},
			0x1d: {K_CTRL, K_5},
			0x1e: {K_CTRL, K_6},
			0x1f: {K_CTRL, K_7},
			0x20: {K_SPACE},
			0x21: {K_SHIFT, K_1},          // !
			0x22: {K_SHIFT, K_APOSTROPHE}, // "
			0x23: {K_SHIFT, K_3},          // #
			0x24: {K_SHIFT, K_4},          // $
			0x25: {K_SHIFT, K_5},          // %
			0x26: {K_SHIFT, K_7},          // &
			0x27: {K_APOSTROPHE},          // '
			0x28: {K_SHIFT, K_9},          // (
			0x29: {K_SHIFT, K_0},          // )
			0x2a: {K_SHIFT, K_8},          // *
			0x2b: {K_SHIFT, K_EQUAL},      // +
			0x2c: {K_COMMA},               // ,
			0x2d: {K_MINUS},               // -
			0x2e: {K_DOT},                 // .
			0x2f: {K_SLASH},               // /
			0x30: {K_0},
			0x31: {K_1},
			0x32: {K_2},
			0x33: {K_3},
			0x34: {K_4},
			0x35: {K_5},
			0x36: {K_6},
			0x37: {K_7},
			0x38: {K_8},
			0x39: {K_9},
			0x3a: {K_SHIFT, K_SEMICOLON}, // :
			0x3b: {K_SEMICOLON},          // ;
			0x3c: {K_SHIFT, K_COMMA},     // <
			0x3d: {K_EQUAL},              // =
			0x3e: {K_SHIFT, K_DOT},       // >
			0x3f: {K_SHIFT, K_SLASH},     // ?
			0x40: {K_SHIFT, K_2},         // @
			0x41: {K_SHIFT, K_A},
			0x42: {K_SHIFT, K_B},
			0x43: {K_SHIFT, K_C},
			0x44: {K_SHIFT, K_D},
			0x45: {K_SHIFT, K_E},
			0x46: {K_SHIFT, K_F},
			0x47: {K_SHIFT, K_G},
			0x48: {K_SHIFT, K_H},
			0x49: {K_SHIFT, K_I},
			0x4a: {K_SHIFT, K_J},
			0x4b: {K_SHIFT, K_K},
			0x4c: {K_SHIFT, K_L},
			0x4d: {K_SHIFT, K_M},
			0x4e: {K_SHIFT, K_N},
			0x4f: {K_SHIFT, K_O},
			0x50: {K_SHIFT, K_P},
			0x51: {K_SHIFT, K_Q},
			0x52: {K_SHIFT, K_R},
			0x53: {K_SHIFT, K_S},
			0x54: {K_SHIFT, K_T},
			0x55: {K_SHIFT, K_U},
			0x56: {K_SHIFT, K_V},
			0x57: {K_SHIFT, K_W},
			0x58: {K_SHIFT, K_X},
			0x59: {K_SHIFT, K_Y},
			0x5a: {K_SHIFT, K_Z},
			0x5b: {K_LEFTBRACE},      // [
			0x5c: {K_BACKSLASH},      // \
			0x5d: {K_RIGHTBRACE},     // ]
			0x5e: {K_SHIFT, K_6},     // ^
			0x5f: {K_SHIFT, K_MINUS}, // _
			0x60: {K_GRAVE},          // `
			0x61: {K_A},
			0x62: {K_B},
			0x63: {K_C},
			0x64: {K_D},
			0x65: {K_E},
			0x66: {K_F},
			0x67: {K_G},
			0x68: {K_H},
			0x69: {K_I},
			0x6a: {K_J},
			0x6b: {K_K},
			0x6c: {K_L},
			0x6d: {K_M},
			0x6e: {K_N},
			0x6f: {K_O},
			0x70: {K_P},
			0x71: {K_Q},
			0x72: {K_R},
			0x73: {K_S},
			0x74: {K_T},
			0x75: {K_U},
			0x76: {K_V},
			0x77: {K_W},
			0x78: {K_X},
			0x79: {K_Y},
			0x7a: {K_Z},
			0x7b: {K_SHIFT, K_LEFTBRACE},  // {
			0x7c: {K_SHIFT, K_BACKSLASH},  // |
			0x7d: {K_SHIFT, K_RIGHTBRACE}, // }
			0x7e: {K_SHIFT, K_GRAVE},      // ~
			0x7f: {K_BACKSPACE},
		}[b[0]]
	} else if len(b) == 3 {
		if b[0] == 0x1b && b[1] == 0x5b {
			return map[byte]Keycodes{
				0x41: {K_UP},
				0x42: {K_DOWN},
				0x43: {K_RIGHT},
				0x44: {K_LEFT},
				0x48: {K_HOME},
				0x46: {K_END},
			}[b[2]]
		} else if b[0] == 0x1b && b[1] == 0x4f {
			return map[byte]Keycodes{
				0x50: {K_F1},
				0x51: {K_F2},
				0x52: {K_F3},
				0x53: {K_F4},
			}[b[2]]
		}
	} else if len(b) == 4 {
		if b[0] == 0x1b && b[1] == 0x5b && b[3] == 0x7e {
			return map[byte]Keycodes{
				0x32: {K_INSERT},
				0x33: {K_DELETE},
				0x35: {K_PAGEUP},
				0x36: {K_PAGEDOWN},
			}[b[2]]
		}
	} else if len(b) == 5 {
		if b[0] == 0x1b && b[1] == 0x5b && b[2] == 0x31 && b[4] == 0x7e {
			return map[byte]Keycodes{
				0x35: {K_F5},
				0x37: {K_F6},
				0x38: {K_F7},
				0x39: {K_F8},
			}[b[3]]
		} else if b[0] == 0x1b && b[1] == 0x5b && b[2] == 0x32 && b[4] == 0x7e {
			return map[byte]Keycodes{
				0x30: {K_F9},
				0x31: {K_F10},
				0x33: {K_F11},
				0x34: {K_F12},
			}[b[3]]
		}
	}
	return Keycodes{}
}

var comboMap = map[Keycodes]Keycodes{
	{K_Q}:      {K_CTRL, K_Q},             // special: exit
	{K_K}:      {K_CTRL, K_K},             // ctrl-k
	{K_F1}:     {K_CTRL, K_ALT, K_F1},     // ctrl-alt-F1
	{K_F2}:     {K_CTRL, K_ALT, K_F2},     // ctrl-alt-F2
	{K_F3}:     {K_CTRL, K_ALT, K_F3},     // ctrl-alt-F3
	{K_F4}:     {K_CTRL, K_ALT, K_F4},     // ctrl-alt-F4
	{K_F5}:     {K_CTRL, K_ALT, K_F5},     // ctrl-alt-F5
	{K_F6}:     {K_CTRL, K_ALT, K_F6},     // ctrl-alt-F6
	{K_F7}:     {K_CTRL, K_ALT, K_F7},     // ctrl-alt-F7
	{K_F8}:     {K_CTRL, K_ALT, K_F8},     // ctrl-alt-F8
	{K_F9}:     {K_CTRL, K_ALT, K_F9},     // ctrl-alt-F9
	{K_F10}:    {K_CTRL, K_ALT, K_F10},    // ctrl-alt-F10
	{K_F11}:    {K_CTRL, K_ALT, K_F11},    // ctrl-alt-F11
	{K_F12}:    {K_CTRL, K_ALT, K_F12},    // ctrl-alt-F12
	{K_DELETE}: {K_CTRL, K_ALT, K_DELETE}, // ctrl-alt-delete
}
var ComboKeycodes = Keycodes{K_CTRL, K_K}
var ComboKeycodesExit = Keycodes{K_CTRL, K_Q}

func CalcCombo(ks Keycodes) Keycodes {
	return comboMap[ks]
}
