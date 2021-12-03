package driver

func EncodeForCH9329(ks Keycodes) []byte {
	var res [14]byte = [14]byte{0x57, 0xab, 0x00, 0x02, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c}
	if ks == EmptyKeycodes {
		return res[:]
	}

	pos := 7
	for _, k := range ks {
		if k == K_L_CTRL {
			res[5] |= 0x01
		} else if k == K_L_SHIFT {
			res[5] |= 0x02
		} else if k == K_L_ALT {
			res[5] |= 0x04
		} else {
			res[pos] = byte(k)
			res[13] += res[pos]
			pos++
		}
	}
	res[13] += res[5]
	return res[:]
}

func EncodeForKCOM3(ks Keycodes) []byte {
	var res [11]byte = [11]byte{0x57, 0xab, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if ks == EmptyKeycodes {
		return res[:]
	}

	pos := 5
	for _, k := range ks {
		if k == K_L_CTRL {
			res[3] |= 0x01
		} else if k == K_L_SHIFT {
			res[3] |= 0x02
		} else if k == K_L_ALT {
			res[3] |= 0x04
		} else {
			res[pos] = byte(k)
			pos++
		}
	}
	return res[:]
}
