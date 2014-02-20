package main

import (
	"github.com/sauerbraten/extinfo"
	"log"
	"net"
	"strconv"
	"time"
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
		case extinfo.BASIC_INFO:
			sendBasicInfo(conn, raddr)
		case extinfo.EXTENDED_INFO:
			extReqType := p.getInt32()
			switch extReqType {
			case extinfo.EXTENDED_INFO_UPTIME:
				sendUptime(conn, raddr)
			case extinfo.EXTENDED_INFO_PLAYER_STATS:
				sendPlayerStats(int(p.getInt32()), conn, raddr)
			case extinfo.EXTENDED_INFO_TEAMS_SCORES:
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
}

func sendUptime(conn *net.UDPConn, raddr *net.UDPAddr) {
	p := NewPacket(0, extinfo.EXTENDED_INFO_UPTIME, extinfo.EXTENDED_INFO_ACK, extinfo.EXTENDED_INFO_VERSION, int(time.Since(state.UpSince)/time.Second))
	conn.WriteToUDP(p.buf, raddr)
}

func sendPlayerStats(cn int, conn *net.UDPConn, raddr *net.UDPAddr) {
	p := &Packet{}
	p.put(0, extinfo.EXTENDED_INFO_PLAYER_STATS, cn, extinfo.EXTENDED_INFO_ACK, extinfo.EXTENDED_INFO_VERSION)

	if cn < -1 || cn > len(clients) {
		p.put(extinfo.EXTENDED_INFO_ERROR)

		n, err := conn.WriteToUDP(p.buf, raddr)
		if err != nil {
			log.Println(err)
		}

		if n != p.len() {
			log.Println("packet length and sent length didn't match!", p.buf)
		}

		return
	}

	p.put(extinfo.EXTENDED_INFO_NO_ERROR)

	n, err := conn.WriteToUDP(p.buf, raddr)
	if err != nil {
		log.Println(err)
	}

	if n != p.len() {
		log.Println("packet length and sent length didn't match!", p.buf)
	}

	p.clear()

	p.put(extinfo.EXTENDED_INFO_PLAYER_STATS_RESPONSE_IDS)

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
			p.put(extinfo.EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS, client.CN, client.Ping, client.Name, client.Team, client.GameState.Frags, client.GameState.Flags, client.GameState.Deaths, client.GameState.Teamkills, client.GameState.Damage*100/max(client.GameState.ShotDamage, 1), client.GameState.Health, client.GameState.Armour, client.GameState.SelectedWeapon, client.Privilege, client.GameState.State)
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
		p.put(extinfo.EXTENDED_INFO_PLAYER_STATS_RESPONSE_STATS, client.CN, client.Ping, client.Name, client.Team, client.GameState.Frags, client.GameState.Flags, client.GameState.Deaths, client.GameState.Teamkills, client.GameState.Damage*100/max(client.GameState.ShotDamage, 1), client.GameState.Health, client.GameState.Armour, client.GameState.SelectedWeapon, client.Privilege, client.GameState.State)
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
