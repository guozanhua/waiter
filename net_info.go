package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

// Protocol constants
const (
	// Constants describing the type of information to query for
	EXTENDED_INFO = 0
	BASIC_INFO    = 1

	NET_INFO_VERSION = 105

	EXTENDED_INFO_ACK      = -1
	EXTENDED_INFO_ERROR    = 1
	EXTENDED_INFO_NO_ERROR = 0

	// Constants describing the type of extended information to query for
	EXTENDED_INFO_UPTIME       = 0
	EXTENDED_INFO_PLAYER_STATS = 1
	EXTENDED_INFO_TEAMS_SCORES = 2

	EXTENDED_INFO_PLAYER_STATS_RESPONSE_IDS   = -10
	EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS = -11
)

func serveStateInfo() {
	// listen for incoming traffic
	laddr, err := net.ResolveUDPAddr("udp", config.ListenAddress+":"+strconv.Itoa(config.ListenPort+1))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("listening for info requests on", laddr.String())

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		buf := make([]byte, 16)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		if n > 5 {
			log.Println("malformed info request:", buf)
			continue
		}

		// process requests

		p := Packet{
			buf: buf,
			pos: 0,
		}

		reqType := p.getInt32()

		switch reqType {
		case BASIC_INFO:
			sendBasicInfo(conn, raddr)
		case EXTENDED_INFO:
			extReqType := p.getInt32()
			switch extReqType {
			case EXTENDED_INFO_UPTIME:
				sendUptime(conn, raddr)
			case EXTENDED_INFO_PLAYER_STATS:
				sendPlayerStats(int(p.getInt32()), conn, raddr)
			case EXTENDED_INFO_TEAMS_SCORES:
				// TODO
			default:
				log.Println("erroneous extinfo type queried:", reqType)
			}
		}
	}
}

func sendBasicInfo(conn *net.UDPConn, raddr *net.UDPAddr) {
	log.Println(conn.RemoteAddr())
	log.Println("basic info requested by", raddr.String())

	p := NewPacket(BASIC_INFO)

	p.put(clients.numberOfClientsInUse())
	p.putInt32(5) // this implementation never sends information about the server being paused or not and the gamespeed
	p.put(NET_INFO_VERSION)
	p.put(state.GameMode)
	p.putInt32(state.TimeLeft / 1000)
	p.put(config.MaxClients)
	p.put(state.MasterMode)
	p.put(state.Map)
	p.put(config.ServerDescription)

	n, err := conn.WriteToUDP(p.buf, raddr)
	if err != nil {
		log.Println(err)
	}

	if n != p.len() {
		log.Println("packet length and sent length didn't match!", p.buf)
	}
}

func sendUptime(conn *net.UDPConn, raddr *net.UDPAddr) {
	p := NewPacket(0, EXTENDED_INFO_UPTIME, EXTENDED_INFO_ACK, NET_INFO_VERSION, int(time.Since(state.UpSince)/time.Second))
	conn.WriteToUDP(p.buf, raddr)
}

func sendPlayerStats(cn int, conn *net.UDPConn, raddr *net.UDPAddr) {
	p := NewPacket(0, EXTENDED_INFO_PLAYER_STATS, cn, EXTENDED_INFO_ACK, NET_INFO_VERSION)

	if cn < -1 || cn > len(clients) {
		p.put(EXTENDED_INFO_ERROR)

		n, err := conn.WriteToUDP(p.buf, raddr)
		if err != nil {
			log.Println(err)
		}

		if n != p.len() {
			log.Println("packet length and sent length didn't match!", p.buf)
		}

		return
	}

	p.put(EXTENDED_INFO_NO_ERROR)

	n, err := conn.WriteToUDP(p.buf, raddr)
	if err != nil {
		log.Println(err)
	}

	if n != p.len() {
		log.Println("packet length and sent length didn't match!", p.buf)
	}

	p.clear()

	p.put(EXTENDED_INFO_PLAYER_STATS_RESPONSE_IDS)

	if cn == -1 {
		for _, client := range clients {
			if !client.Joined {
				continue
			}
			p.put(client.CN)
		}
	} else {
		p.put(clients[ClientNumber(cn)].CN)
	}

	n, err = conn.WriteToUDP(p.buf, raddr)
	if err != nil {
		log.Println(err)
	}

	if n != p.len() {
		log.Println("packet length and sent length didn't match!", p.buf)
	}

	p.clear()

	if cn == -1 {
		for _, client := range clients {
			if !client.Joined {
				continue
			}
			p.put(EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS, client.CN, client.Ping, client.Name, client.Team, client.GameState.Frags, client.GameState.Flags, client.GameState.Deaths, client.GameState.Teamkills, client.GameState.Damage*100/max(client.GameState.ShotDamage, 1), client.GameState.Health, client.GameState.Armour, client.GameState.SelectedWeapon, client.Privilege, client.GameState.State)
			if config.SendClientIPsViaExtinfo {
				p.put(client.Peer.Address.IP[:2]) // only 3 first bytes
			} else {
				p.put(0, 0, 0) // 3 times 0x0
			}

			n, err = conn.WriteToUDP(p.buf, raddr)
			if err != nil {
				log.Println(err)
			}

			if n != p.len() {
				log.Println("packet length and sent length didn't match!", p.buf)
			}

			p.clear()
		}
	} else {
		client := clients[ClientNumber(cn)]
		p.put(EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS, client.CN, client.Ping, client.Name, client.Team, client.GameState.Frags, client.GameState.Flags, client.GameState.Deaths, client.GameState.Teamkills, client.GameState.Damage*100/max(client.GameState.ShotDamage, 1), client.GameState.Health, client.GameState.Armour, client.GameState.SelectedWeapon, client.Privilege, client.GameState.State)
		if config.SendClientIPsViaExtinfo {
			p.put(client.Peer.Address.IP[:2]) // only 3 first bytes
		} else {
			p.put(0, 0, 0) // 3 times 0x0
		}

		n, err = conn.WriteToUDP(p.buf, raddr)
		if err != nil {
			log.Println(err)
		}

		if n != p.len() {
			log.Println("packet length and sent length didn't match!", p.buf)
		}

		p.clear()
	}
}

func max(i, j int32) int32 {
	if i > j {
		return i
	} else {
		return j
	}
}
