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

func VT100Decode(b []byte) Keycodes {
	if len(b) == 1 {
		if 'a' <= b[0] && b[0] <= 'z' {
			return Keycodes{b[0] - 'a' + K_A}
		} else if 'A' <= b[0] && b[0] <= 'Z' {
			return Keycodes{K_SHIFT, b[0] - 'A' + K_A}
		} else if 0x01 <= b[0] && b[0] <= 0x1a {
			return Keycodes{K_CTRL, b[0] - 0x01 + K_A}
		} else if '0' <= b[0] && b[0] <= '9' {
			return Keycodes{b[0] - '0' + K_0}
		} else if 0x1b <= b[0] && b[0] <= 0x1f {
			return Keycodes{K_CTRL, b[0] - 0x1b + K_3}
		} else {
			return map[byte]Keycodes{
				0x1b: {K_ESC},
				'`':  {K_GRAVE},
				'~':  {K_SHIFT, K_GRAVE},
				'!':  {K_SHIFT, K_1},
				'@':  {K_SHIFT, K_2},
				'#':  {K_SHIFT, K_3},
				'$':  {K_SHIFT, K_4},
				'%':  {K_SHIFT, K_5},
				'^':  {K_SHIFT, K_6},
				'&':  {K_SHIFT, K_7},
				'*':  {K_SHIFT, K_8},
				'(':  {K_SHIFT, K_9},
				')':  {K_SHIFT, K_0},
				'-':  {K_MINUS},
				'_':  {K_SHIFT, K_MINUS},
				'=':  {K_EQUAL},
				'+':  {K_SHIFT, K_EQUAL},
				0x7f: {K_BACKSPACE},
				'\t': {K_TAB},
				'[':  {K_LEFTBRACE},
				'{':  {K_SHIFT, K_LEFTBRACE},
				']':  {K_RIGHTBRACE},
				'}':  {K_SHIFT, K_RIGHTBRACE},
				'\\': {K_BACKSLASH},
				'|':  {K_SHIFT, K_BACKSLASH},
				';':  {K_SEMICOLON},
				':':  {K_SHIFT, K_SEMICOLON},
				'\'': {K_APOSTROPHE},
				'"':  {K_SHIFT, K_APOSTROPHE},
				'\r': {K_ENTER},
				',':  {K_COMMA},
				'<':  {K_SHIFT, K_COMMA},
				'.':  {K_DOT},
				'>':  {K_SHIFT, K_DOT},
				'/':  {K_SLASH},
				'?':  {K_SHIFT, K_SLASH},
				' ':  {K_SPACE},
				0x00: {K_SHIFT, K_2},
			}[b[0]]
		}
	} else if len(b) == 3 {
		if b[0] == 0x1b && b[1] == 0x5b {
			return map[byte]Keycodes{
				'A': {K_UP},
				'B': {K_DOWN},
				'C': {K_RIGHT},
				'D': {K_LEFT},
				'H': {K_HOME},
				'F': {K_END},
			}[b[2]]
		} else if b[0] == 0x1b && b[1] == 0x4f {
			return map[byte]Keycodes{
				'P': {K_F1},
				'Q': {K_F2},
				'R': {K_F3},
				'S': {K_F4},
			}[b[2]]
		}
	} else if len(b) == 4 {
		if b[0] == 0x1b && b[1] == 0x5b && b[3] == 0x7e {
			return map[byte]Keycodes{
				'2': {K_INSERT},
				'3': {K_DELETE},
				'5': {K_PAGEUP},
				'6': {K_PAGEDOWN},
			}[b[2]]
		}
	} else if len(b) == 5 {
		if b[0] == 0x1b && b[1] == 0x5b && b[2] == 0x31 && b[4] == 0x7e {
			return map[byte]Keycodes{
				'5': {K_F5},
				'7': {K_F6},
				'8': {K_F7},
				'9': {K_F8},
			}[b[3]]
		} else if b[0] == 0x1b && b[1] == 0x5b && b[2] == 0x32 && b[4] == 0x7e {
			return map[byte]Keycodes{
				'0': {K_F9},
				'1': {K_F10},
				'3': {K_F11},
				'4': {K_F12},
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

func CalcCombo(ks Keycodes) Keycodes {
	return comboMap[ks]
}
