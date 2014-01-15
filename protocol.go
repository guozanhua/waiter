package main

import (
	"log"
)

const PROTOCOL_VERSION int32 = 259

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

		case Packet:
			p = append(p, v...)

		case PlayerPosition:
			p = append(p, Packet(v)...)
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
		nmc := NetworkMessageCode(p.getInt32())

		if !isValidNetworkMessageCode(nmc, client) {
			client.disconnect(DISC_MSGERR)
			return
		}

		switch nmc {
		case N_JOIN:
			// client sends intro and wants to join the game
			log.Println("received N_JOIN")
			client.join(p.getString(), p.getInt32(), p.getString(), p.getString(), p.getString())
			break outer

		case N_AUTHANS:
			// client sends answer to auth challenge
			log.Println("received N_AUTHANS")
			break outer

		case N_PING:
			// client pinging server
			log.Println("received N_PING")
			p.getInt32()
			break outer

		case N_POS:
			// client sending his position in the world
			//log.Println("received N_POS")
			log.Println("pos of", client.CN, N_POS, p)
			client.GameState.Position = PlayerPosition(append([]byte{byte(N_POS)}, p...))
			log.Println("pos of", client.CN, "set to", client.GameState.Position)
			break outer

		default:
			log.Println(nmc, "on channel", channelId, p)
			break outer
		}
	}

	return
}
