package driver

func EncodeForCH9329_K(b []byte) []byte {
	return b
}

func EncodeForCH9329_M(b []byte) []byte {
	var offset = 30
	var res [11]byte = [11]byte{0x57, 0xab, 0x00, 0x05, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x0d}

	if len(b) == 3 && b[0] == 0x1b && b[1] == 0x5b {
		switch b[2] {
		case 0x41: // up
			res[8] = byte(-offset & 0xff)
			res[10] = byte((res[8] + res[10]) & 0xff)
			return res[:]
		case 0x42: // down
			res[8] = byte(offset & 0xff)
			res[10] = byte((res[8] + res[10]) & 0xff)
			return res[:]
		case 0x43: // right
			res[7] = byte(offset & 0xff)
			res[10] = byte((res[7] + res[10]) & 0xff)
			return res[:]
		case 0x44: // left
			res[7] = byte(-offset & 0xff)
			res[10] = byte((res[7] + res[10]) & 0xff)
			return res[:]
		}
	}
	return nil
}
