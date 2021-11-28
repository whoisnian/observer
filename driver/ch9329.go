package driver

// 模拟按下“A”键: 0x57、0xAB、0x00、0x02、0x08、0x00、0x00、0x04、0x00、
// 0x00、0x00、0x00、0x00、0x10。
// 模拟释放“A”键: 0x57、0xAB、0x00、0x02、0x08、0x00、0x00、0x00、0x00、
// 0x00、0x00、0x00、0x00、0x0C。
// BIT7  BIT6  BIT5    BIT4   BIT3  BIT2  BIT1    BIT0
// R_win R_alt R_shift R_ctrl L_win L_alt L_shift L_ctrl

func EncodeForCH9329_K(b []byte) []byte {
	var res [14]byte = [14]byte{0x57, 0xab, 0x00, 0x02, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c}
	if b == nil {
		return res[:]
	} else if len(b) == 1 {
		if 'a' <= b[0] && b[0] <= 'z' {
			res[7] = b[0] - 'a' + 0x04
		} else if 'A' <= b[0] && b[0] <= 'Z' {
			res[7] = b[0] - 'A' + 0x04
			res[5] = 0x02 // L_shift
		} else if 0x01 <= b[0] && b[0] <= 0x1a {
			res[7] = b[0] - 0x01 + 0x04
			res[5] = 0x01 // L_ctrl
		} else if '1' <= b[0] && b[0] <= '9' {
			res[7] = b[0] - '1' + 0x1e
		} else if 0x1b <= b[0] && b[0] <= 0x1f {
			res[7] = b[0] - 0x1b + 0x20
			res[5] = 0x01 // L_ctrl
		} else {
			switch b[0] {
			case '0':
				res[7] = 0x27
			case ')':
				res[7] = 0x27
				res[5] = 0x02 // L_shift
			case '!':
				res[7] = 0x1e
				res[5] = 0x02 // L_shift
			case '@':
				res[7] = 0x1f
				res[5] = 0x02 // L_shift
			case '#':
				res[7] = 0x20
				res[5] = 0x02 // L_shift
			case '$':
				res[7] = 0x21
				res[5] = 0x02 // L_shift
			case '%':
				res[7] = 0x22
				res[5] = 0x02 // L_shift
			case '^':
				res[7] = 0x23
				res[5] = 0x02 // L_shift
			case '&':
				res[7] = 0x24
				res[5] = 0x02 // L_shift
			case '*':
				res[7] = 0x25
				res[5] = 0x02 // L_shift
			case '(':
				res[7] = 0x26
				res[5] = 0x02 // L_shift
			case '-':
				res[7] = 0x2d
			case '_':
				res[7] = 0x2d
				res[5] = 0x02 // L_shift
			case '=':
				res[7] = 0x2e
			case '+':
				res[7] = 0x2e
				res[5] = 0x02 // L_shift
			case '`':
				res[7] = 0x35
			case '~':
				res[7] = 0x35
				res[5] = 0x02 // L_shift
			case 0x7f:
				res[7] = 0x2a
			case '[':
				res[7] = 0x2f
			case '{':
				res[7] = 0x2f
				res[5] = 0x02 // L_shift
			case ']':
				res[7] = 0x30
			case '}':
				res[7] = 0x30
				res[5] = 0x02 // L_shift
			case '\\':
				res[7] = 0x31
			case '|':
				res[7] = 0x31
				res[5] = 0x02 // L_shift
			case ';':
				res[7] = 0x33
			case ':':
				res[7] = 0x33
				res[5] = 0x02 // L_shift
			case '\'':
				res[7] = 0x34
			case '"':
				res[7] = 0x34
				res[5] = 0x02 // L_shift
			case '\r':
				res[7] = 0x28
			case ',':
				res[7] = 0x36
			case '<':
				res[7] = 0x36
				res[5] = 0x02 // L_shift
			case '.':
				res[7] = 0x37
			case '>':
				res[7] = 0x37
				res[5] = 0x02 // L_shift
			case '/':
				res[7] = 0x38
			case '?':
				res[7] = 0x38
				res[5] = 0x02 // L_shift
			case ' ':
				res[7] = 0x2c
			case 0x1b:
				res[7] = 0x29
			case '\t':
				res[7] = 0x29
			case 0x00:
				res[7] = 0x1f
				res[5] = 0x01 // L_ctrl
			}
		}
	} else if len(b) == 2 {
	} else if len(b) == 3 {
		if b[0] == 0x1b && b[1] == 0x5b {
			switch b[2] {
			case 0x41: // up
				res[7] = 0x52
			case 0x42: // down
				res[7] = 0x51
			case 0x43: // right
				res[7] = 0x4f
			case 0x44: // left
				res[7] = 0x50
			case 0x48: // home
				res[7] = 0x4a
			case 0x46: // end
				res[7] = 0x4d
			}
		} else if b[0] == 0x1b && b[1] == 0x4f {
			switch b[2] {
			case 0x50: // f1
				res[7] = 0x3a
			case 0x51: // f2
				res[7] = 0x3b
			case 0x52: // f3
				res[7] = 0x3c
			case 0x53: // f4
				res[7] = 0x3d
			}
		}
	} else if len(b) == 4 {
		if b[0] == 0x1b && b[1] == 0x5b && b[3] == 0x7e {
			switch b[2] {
			case 0x32: // insert
				res[7] = 0x49
			case 0x33: // delete
				res[7] = 0x4c
			case 0x35: // page up
				res[7] = 0x4b
			case 0x36: // page down
				res[7] = 0x4e
			}
		}
	} else if len(b) == 5 {
		if b[0] == 0x1b && b[1] == 0x5b && b[2] == 0x31 && b[4] == 0x7e {
			switch b[3] {
			case 0x35: // f5
				res[7] = 0x3e
			case 0x37: // f6
				res[7] = 0x3f
			case 0x38: // f7
				res[7] = 0x40
			case 0x39: // f8
				res[7] = 0x41
			}
		} else if b[0] == 0x1b && b[1] == 0x5b && b[2] == 0x32 && b[4] == 0x7e {
			switch b[3] {
			case 0x30: // f9
				res[7] = 0x42
			case 0x31: // f10
				res[7] = 0x43
			case 0x33: // f11
				res[7] = 0x44
			case 0x34: // f12
				res[7] = 0x45
			}
		}
	}
	res[13] += res[5] + res[7]
	return res[:]
}

func EncodeForCH9329_M(b []byte) []byte {
	var offset byte = 30
	var res [11]byte = [11]byte{0x57, 0xab, 0x00, 0x05, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x0d}

	if len(b) == 3 && b[0] == 0x1b && b[1] == 0x5b {
		switch b[2] {
		case 0x41: // up
			res[8] = -offset
		case 0x42: // down
			res[8] = offset
		case 0x43: // right
			res[7] = offset
		case 0x44: // left
			res[7] = -offset
		}
		res[10] += res[7] + res[8]
		return res[:]
	}
	return nil
}
