package main

import (
	"log"
)

var clientPositions map[ClientNumber]*PlayerPosition

func sendWorldState() {
	clientPositions = map[ClientNumber]*PlayerPosition{}
	clientsInUse := 0

	for _, client := range clients {
		if !client.Joined || len(client.GameState.Position) == 0 {
			continue
		}

		log.Println("adding", client.CN, "to world state")

		clientPositions[client.CN] = &client.GameState.Position
		clientsInUse++
	}

	if clientsInUse == 0 {
		return
	}

	log.Println("world state", clientPositions)

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
		sendf(client, false, 0, positions...)
	}
}
