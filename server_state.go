package main

import (
	"./enet"
)

type ServerState struct {
	MasterMode  MasterMode
	GameMode    GameMode
	Map         string
	TimeLeft    int32 // in milliseconds
	NotGotItems bool
	HasMaster   bool
}

func (state *ServerState) changeMap(mapName string) {
	state.NotGotItems = true
	state.Map = mapName
	clients.send(enet.PACKET_FLAG_RELIABLE, 1, N_MAPCHANGE, state.Map, state.GameMode, state.NotGotItems)
	clients.send(enet.PACKET_FLAG_RELIABLE, 1, N_TIMELEFT, state.TimeLeft/1000)
	for _, c := range clients {
		if !c.InUse || c.GameState.State == CS_SPECTATOR {
			continue
		}
		c.GameState.reset()
		c.sendSpawnState()
	}
}
