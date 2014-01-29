package main

import (
	"log"
)

const PROTOCOL_VERSION int32 = 259

// makes a new packet containing all the values in args
func makePacket(args ...interface{}) (p Packet) {
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

		case WeaponNumber:
			p.putInt32(int32(v))

		case ArmourType:
			p.putInt32(int32(v))

		case DisconnectReason:
			p.putInt32(int32(v))

		case Packet:
			p.putBytes(v.buf)

		case PlayerPosition:
			p.putBytes(v.buf)

		default:
			log.Printf("unhandled type %T of arg %v\n", v, v)
		}
	}

	return
}

// parses a packet and decides what to do based on the network message code at the front of the packet
func parsePacket(fromCN ClientNumber, channelId uint8, p Packet) {
	if fromCN < 0 || channelId > 2 {
		return
	}

	client := clients[fromCN]

outer:
	for p.pos < p.len() {
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
			//break outer

		case N_AUTHANS:
			// client sends answer to auth challenge
			log.Println("received N_AUTHANS")
			//break outer

		case N_PING:
			// client pinging server
			//log.Println("received N_PING")
			p.getInt32()
			//break outer

		case N_POS:
			// client sending his position in the world
			client.GameState.Position = PlayerPosition(p)
			break outer

		case N_TEXT:
			// client sending chat message → add to client's buffered messages
			client.GameState.BufferedPackets = append(client.GameState.BufferedPackets, makePacket(N_TEXT, p.getString()))

		case N_SAYTEAM:
			// client sending team chat message → pass on to team immediatly
			client.sendToTeam(true, 1, N_SAYTEAM, client.CN, p.getString())

		case N_WEAPONSELECT:
			// player changing weapon
			selectedWeapon := WeaponNumber(p.getInt32())
			if client.GameState.State != CS_ALIVE {
				break
			}

			if selectedWeapon >= WPN_SAW && selectedWeapon <= WPN_PISTOL {
				client.GameState.SelectedWeapon = selectedWeapon
			} else {
				client.GameState.SelectedWeapon = WPN_PISTOL
			}

			client.GameState.BufferedPackets = append(client.GameState.BufferedPackets, makePacket(N_WEAPONSELECT, selectedWeapon))

		default:
			log.Println(p, "on channel", channelId)
			break outer
		}
	}

	return
}
