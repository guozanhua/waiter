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

		case GameState:
			p.putInt32(v.LifeSequence)
			p.putInt32(v.Health)
			p.putInt32(v.MaxHealth)
			p.putInt32(v.Armour)
			p.putInt32(int32(v.ArmourType))
			p.putInt32(int32(v.SelectedWeapon))
			for _, ammo := range v.Ammo {
				p.putInt32(ammo)
			}

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
			// client pinging server → send pong
			client.send(false, 1, N_PONG, p.getInt32())

		case N_CLIENTPING:
			// client sending the amount of LAG he measured to the server → broadcast by adding to client's buffered messages
			client.Ping = p.getInt32()
			client.GameState.BufferedPackets = append(client.GameState.BufferedPackets, makePacket(N_CLIENTPING, client.Ping))

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

		case N_MAPCRC:
			// client sends crc hash of his map file
			//clientMapName := p.getString()
			//clientMapCRC := p.getInt32()
			p.getString()
			p.getInt32()

		case N_SPAWN:
			client.spawn(p.getInt32(), p.getInt32())
			client.GameState.BufferedPackets = append(client.GameState.BufferedPackets, makePacket(N_SPAWN, client.GameState))

		case N_WEAPONSELECT:
			// player changing weapon
			selectedWeapon := WeaponNumber(p.getInt32())
			client.GameState.selectWeapon(selectedWeapon)

			// add to client's message buffer to broadcast to other clients
			client.GameState.BufferedPackets = append(client.GameState.BufferedPackets, makePacket(N_WEAPONSELECT, selectedWeapon))

		default:
			log.Println(p, "on channel", channelId)
			break outer
		}
	}

	return
}
