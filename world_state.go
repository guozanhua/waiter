package main

import (
	"log"
)

var clientPositions map[ClientNumber]*PlayerPosition

func sendPositions() {
	clientPositions = map[ClientNumber]*PlayerPosition{}
	clientsInUse := 0

	for _, client := range clients {
		if !client.Joined || len(client.GameState.Position.buf) == 0 {
			continue
		}

		clientPositions[client.CN] = &client.GameState.Position
		clientsInUse++
	}

	if clientsInUse == 0 {
		return
	}

	for _, client := range clients {
		if !client.Joined {
			continue
		}

		positions := []interface{}{}

		for otherCN, positionPacket := range clientPositions {
			if client.CN == otherCN {
				continue
			}

			positions = append(positions, *positionPacket)
		}

		// send on channel 0
		client.send(false, 0, positions...)
	}
}

var clientPackets map[ClientNumber][]Packet

func sendNetworkMessages() {
	clientPackets = map[ClientNumber][]Packet{}
	clientsInUse := 0
	reliablePacketPresent := false

	for _, client := range clients {
		if !client.Joined || len(client.GameState.BufferedPackets) == 0 {
			continue
		}

		log.Println("adding packets by", client.CN)

		clientPackets[client.CN] = client.GameState.BufferedPackets
		client.GameState.BufferedPackets = []Packet{}
		reliablePacketPresent = client.GameState.HasReliablePacket
		clientsInUse++
	}

	if clientsInUse == 0 {
		return
	}

	for _, client := range clients {
		if !client.Joined {
			continue
		}

		packets := []interface{}{}

		for otherCN, packetBuffer := range clientPackets {
			if client.CN == otherCN {
				continue
			}

			packets = append(packets, N_CLIENT, otherCN)

			for _, packet := range packetBuffer {
				packets = append(packets, packet)
			}
		}

		// send on channel 1
		client.send(reliablePacketPresent, 1, packets...)
	}
}
