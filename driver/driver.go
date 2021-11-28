package driver

type Key uint32
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

	K_SHIFT
	K_CTRL
	K_ALT
	K_SPACE
	K_TAB
	K_ENTER
	K_BACKSPACE

	K_INSERT
	K_DELETE
	K_HOME
	K_END
	K_PAGEUP
	K_PAGEDOWN

	K_UP
	K_DOWN
	K_LEFT
	K_RIGHT
)

func VT100Decode(b []byte) Keycodes {
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

func CalcCombo(b []byte) Keycodes {
	if code, ok := comboMap[VT100Decode(b)]; ok {
		return code
	} else {
		return Keycodes{}
	}
}
