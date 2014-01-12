package main

import (
	"log"
)

const PROTOCOL_VERSION int32 = 259

type Packet []byte

// encodes an int32 and appends it to the packet
func (p *Packet) putInt32(i int32) {
	if i < 128 && i > -127 {
		*p = append(*p, byte(i))
	} else if i < 0x8000 && i >= -0x8000 {
		*p = append(*p, 0x80, byte(i), byte(i>>8))
	} else {
		*p = append(*p, 0x81, byte(i), byte(i>>8), byte(i>>16), byte(i>>24))
	}
}

// appends a string to the packet
func (p *Packet) putString(s string) {
	for _, c := range s {
		p.putInt32(int32(c))
	}
	p.putInt32(0)
}

// returns the first byte in the Packet and shortens the packet slice to no longer include this byte
func (p *Packet) popByte() byte {
	b := (*p)[0]
	*p = (*p)[1:]
	return b
}

// decodes an int32 and removes it from the front of the packet
func (p *Packet) popInt32() int32 {
	i := int32(p.popByte())

	if i == 0x80 {
		return int32(p.popByte()) + (int32(p.popByte()) << 8)
	} else if i == 0x81 {
		return int32(p.popByte()) + (int32(p.popByte()) << 8) + (int32(p.popByte()) << 16) + (int32(p.popByte()) << 24)
	} else {
		return i
	}
}

// decodes an int32 and removes it from the front of the packet, using the different compression meant for uint32s
func (p *Packet) popUint32() int32 {
	i := int32(p.popByte())
	if i >= 0x80 {
		i += int32(p.popByte()<<7) - 0x80
		if i >= (1 << 14) {
			i += int32(p.popByte()<<14) - (1 << 14)
		}
		if i >= (1 << 21) {
			i += int32(p.popByte()<<21) - (1 << 21)
		}
		if i >= (1 << 28) {
			i |= -(1 << 28)
		}
	}

	return i
}

// reads a string and removes it from the front of the packet
func (p *Packet) popString() string {
	buf := []byte{}

	for b := p.popByte(); b != 0x00; b = p.popByte() {
		buf = append(buf, b)
	}

	return string(buf)
}

// sends a (reliable if desired) packet to client, containing all the parts
func sendf(client *Client, reliable bool, channel uint8, args ...interface{}) {
	p := Packet{}

	for _, arg := range args {
		switch v := arg.(type) {
		case int32:
			p.putInt32(v)

		case int:
			p.putInt32(int32(v))

		case byte:
			p.putInt32(int32(v))

		case bool:
			if v {
				p.putInt32(1)
			} else {
				p.putInt32(0)
			}

		case string:
			p.putString(v)

		case NetworkMessageCode:
			p.putInt32(int32(v))

		case MasterMode:
			p.putInt32(int32(v))

		case GameMode:
			p.putInt32(int32(v))

		case ClientNumber:
			p.putInt32(int32(v))

		case ClientState:
			p.putInt32(int32(v))

		case GunNumber:
			p.putInt32(int32(v))

		case ArmourType:
			p.putInt32(int32(v))

		case DisconnectReason:
			p.putInt32(int32(v))
		}
	}

	client.send(&p, reliable, channel)
}

// parses a packet and decides what to do based on the network message code at the front of the packet
func parsePacket(fromCN ClientNumber, channelId uint8, p Packet) {
	if fromCN < 0 || channelId > 2 {
		return
	}

	client := clients[fromCN]

outer:
	for len(p) > 0 {
		nmc := NetworkMessageCode(p.popInt32())

		if !isValidNetworkMessageCode(nmc, client) {
			client.disconnect(DISC_MSGERR)
			return
		}

		switch nmc {
		case N_JOIN:
			// client sends intro and wants to join the game
			log.Println("received N_JOIN")
			client.tryToJoin(p.popString(), p.popInt32(), p.popString(), p.popString(), p.popString())
			break outer

		case N_AUTHANS:
			// client sends answer to auth challenge
			log.Println("received N_AUTHANS")
			break outer

		case N_PING:
			// client pinging server
			log.Println("received N_PING")
			p.popInt32()
			break outer

		case N_POS:
			// client sending his position in the world
			//log.Println("received N_POS")
			log.Println("pos of", client.CN, N_POS, p)
			client.GameState.Position = append([]byte{byte(N_POS)}, []byte(p)...)
			log.Println("pos of", client.CN, "set to", client.GameState.Position)
			break outer

		default:
			log.Println(nmc, "on channel", channelId, p)
			break outer
		}
	}

	return
}
