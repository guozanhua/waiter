package main

type Packet struct {
	buf []byte
	pos int
}

func (p *Packet) len() int {
	return len(p.buf)
}

// Appends a []byte to the end of the packet.
func (p *Packet) putBytes(b []byte) {
	p.buf = append(p.buf, b...)
}

// Encodes an int32 and appends it to the packet.
func (p *Packet) putInt32(i int32) {
	if i < 128 && i > -127 {
		p.buf = append(p.buf, byte(i))
	} else if i < 0x8000 && i >= -0x8000 {
		p.buf = append(p.buf, 0x80, byte(i), byte(i>>8))
	} else {
		p.buf = append(p.buf, 0x81, byte(i), byte(i>>8), byte(i>>16), byte(i>>24))
	}
}

// Appends a string to the packet.
func (p *Packet) putString(s string) {
	for _, c := range s {
		p.putInt32(int32(c))
	}
	p.putInt32(0)
}

// Returns the first byte in the Packet.
func (p *Packet) getByte() byte {
	b := p.buf[p.pos]
	p.pos++
	return b
}

// Decodes an int32 and increases the position index accordingly.
func (p *Packet) getInt32() int32 {
	i := int32(p.getByte())

	if i == 0x80 {
		return int32(p.getByte()) + (int32(p.getByte()) << 8)
	} else if i == 0x81 {
		return int32(p.getByte()) + (int32(p.getByte()) << 8) + (int32(p.getByte()) << 16) + (int32(p.getByte()) << 24)
	} else {
		return i
	}
}

// Decodes an int32 using the different compression meant for uint32s and increases the position index accordingly.
func (p *Packet) getUint32() int32 {
	i := int32(p.getByte())
	if i >= 0x80 {
		i += int32(p.getByte()<<7) - 0x80
		if i >= (1 << 14) {
			i += int32(p.getByte()<<14) - (1 << 14)
		}
		if i >= (1 << 21) {
			i += int32(p.getByte()<<21) - (1 << 21)
		}
		if i >= (1 << 28) {
			i |= -(1 << 28)
		}
	}

	return i
}

// Reads a string from the packet and increases the position index accordingly.
func (p *Packet) getString() string {
	buf := []byte{}

	for b := p.getByte(); b != 0x00; b = p.getByte() {
		buf = append(buf, b)
	}

	return string(buf)
}
