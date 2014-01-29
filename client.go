package main

import (
	"./enet"
	"log"
)

// Cumulative type for the collection of clients.
type Clients map[ClientNumber]*Client

// Sends a packet to all clients currently in use.
func (cs Clients) send(reliable bool, channel uint8, args ...interface{}) {
	for _, c := range cs {
		if !c.InUse {
			continue
		}
		c.send(reliable, channel, args...)
	}
}

// A player's cn
type ClientNumber int32

// Describes a client's level of privilege.
type ClientPrivilege int32

const (
	PRIV_NONE ClientPrivilege = iota
	PRIV_MASTER
	PRIV_AUTH
	PRIV_ADMIN
)

// Describes a client.
type Client struct {
	CN                  ClientNumber
	Name                string
	Team                string
	PlayerModel         int32
	Privilege           ClientPrivilege
	GameState           GameState
	Joined              bool             // true if the player is actually in the game
	HasToAuthForConnect bool             // true if the server is private or demands auth-on-connect and the client has not yet joined the actual game
	ReasonWhyAuthNeeded DisconnectReason // e.g. server is in private mode
	AI                  bool             // wether this is a bot or not
	AISkill             int32
	InUse               bool // true if this client object is in use in the actual game
	Peer                *enet.Peer
	SessionId           int32
}

// Connects an ENet peer to a client object. If no unused client object can be found, a new one is created and added to the global set of clients.
func addClient(peer *enet.Peer) *Client {
	var client *Client

	// re-use unused client object with low cn
	for _, client = range clients {
		if !client.InUse {
			client.InUse = true
			return client
		}
	}

	client = &Client{
		CN:        ClientNumber(len(clients)),
		InUse:     true,
		Peer:      peer,
		SessionId: rng.Int31(),
		Team:      "good", // TODO: select weaker team
		GameState: GameState{},
	}

	clients[client.CN] = client

	log.Println("added client:", client, "with peer", client.Peer)

	return client
}

// Send a packet to a client (reliable if desired) over the specified channel.
func (client *Client) send(reliable bool, channel uint8, args ...interface{}) {
	p := makePacket(args...)

	mustFlush = true

	var flags enet.PacketFlag
	if reliable {
		flags |= enet.PACKET_FLAG_RELIABLE
	}

	if channel == 1 {
		log.Println(p, "→", client.CN)
	}

	client.Peer.Send(p.buf, flags, channel)
}

// Send a packet to a client's team (reliable if desired) over the specified channel.
func (client *Client) sendToTeam(reliable bool, channel uint8, args ...interface{}) {
	for _, c := range clients {
		if c == client || !c.InUse || c.Team != client.Team {
			continue
		}
		c.send(reliable, channel, args...)
	}
}

// Sends a packet to all clients but the client himself.
func (client *Client) sendToAllOthers(reliable bool, channel uint8, args ...interface{}) {
	for _, c := range clients {
		if c == client || !c.InUse {
			continue
		}
		c.send(reliable, channel, args...)
	}
}

// Sends basic server info to the client.
func (client *Client) sendServerInfo() {
	parts := []interface{}{N_SERVINFO, client.CN, PROTOCOL_VERSION, client.SessionId}

	if config.ServerPassword != "" {
		parts = append(parts, 1)
	} else {
		parts = append(parts, 0)
	}

	parts = append(parts, config.ServerDescription, config.ServerAuthDomains[0])

	client.send(true, 1, parts...)
}

// Sends 'welcome' information to a newly joined client like map, mode, time left, other players, etc.
func (client *Client) sendWelcome() {
	parts := []interface{}{N_WELCOME}

	// send currently played map
	parts = append(parts, N_MAPCHANGE, state.Map, state.GameMode, state.NotGotItems)

	// send time left in this round
	parts = append(parts, N_TIMELEFT, state.TimeLeft/1000)

	// send list of clients which have privilege higher than PRIV_NONE and their respecitve privilege level
	if state.HasMaster {
		parts = append(parts, N_CURRENTMASTER, state.MasterMode)
		for _, c := range clients {
			if c.Privilege > PRIV_NONE {
				parts = append(parts, c.CN, c.Privilege)
			}
		}
		parts = append(parts, -1)
	}

	// tell the client what team he was put in by the server
	parts = append(parts, N_SETTEAM, client.CN, client.Team, -1)

	// tell the client how to spawn (what health, what armour, what weapons, what ammo, etc.)
	if client.GameState.State == CS_SPECTATOR {
		parts = append(parts, N_SPECTATOR, client.CN, 1)
	} else {
		// TODO: handle spawn delay (e.g. in ctf modes)
		parts = append(parts, N_SPAWNSTATE, client.CN, client.GameState.LifeSequence, client.GameState.Health, client.GameState.MaxHealth, client.GameState.Armour, client.GameState.ArmourType, client.GameState.SelectedGun)
		for _, gn := range GunsWithAmmo {
			parts = append(parts, client.GameState.Ammo[gn])
		}
	}

	// send other players' state (frags, flags, etc.)
	parts = append(parts, N_RESUME)
	for _, c := range clients {
		log.Println(c)
		if c == client || !c.InUse {
			continue
		}

		parts = append(parts, c.CN, c.GameState.State, c.GameState.Frags, c.GameState.Flags, c.GameState.QuadTimeLeft, c.GameState.LifeSequence, c.GameState.Health, c.GameState.MaxHealth, c.GameState.Armour, c.GameState.ArmourType, c.GameState.SelectedGun)
		for _, gn := range GunsWithAmmo {
			parts = append(parts, c.GameState.Ammo[gn])
		}
	}
	parts = append(parts, -1)

	// send other client's state (name, team, playermodel)
	for _, c := range clients {
		if c == client || !c.InUse {
			continue
		}
		parts = append(parts, N_INITCLIENT, c.CN, c.Name, c.Team, c.PlayerModel)
	}

	log.Println(parts)

	client.send(true, 1, parts...)
}

// Tries to let a client join the current game, using the data the client provided with his N_JOIN packet.
func (client *Client) join(name string, playerModel int32, hash string, authDomain string, authName string) {
	// TODO: check server password hash

	// check for mandatory connect auth
	if client.HasToAuthForConnect {
		if authDomain != config.ServerAuthDomains[0] {
			// client has no authkey for the server domain
			// TODO: disconnect client with disconnect reason
		}
	}

	// player may join
	client.Joined = true
	client.Name = name
	client.PlayerModel = playerModel

	client.GameState.spawn(state.GameMode)

	if state.MasterMode == MM_LOCKED {
		client.GameState.State = CS_SPECTATOR
	}

	// send welcome packet
	client.sendWelcome()

	// inform other clients that a new client joined
	client.informOthersOfJoin()

	log.Printf("join: %s (%d)\n", name, client.CN)

	return
}

// For when a client disconnects deliberately.
func (client *Client) leave() {
	log.Printf("leave: %s (%d)\n", client.Name, client.CN)
	client.disconnect(DISC_NONE)
}

// Tells other clients that the client disconnected, giving a disconnect reason in case it's not a normal leave.
func (client *Client) disconnect(reason DisconnectReason) {
	if !client.InUse {
		return
	}

	// inform others
	client.informOthersOfDisconnect(reason)

	if reason != DISC_NONE {
		log.Printf("disconnected: %s (%d) - %s", DisconnectReasons[reason])
	}

	client.Peer.Disconnect(uint32(reason))

	client.reset()
}

// Informs all other clients that a client joined the game.
func (client *Client) informOthersOfJoin() {
	client.sendToAllOthers(true, 1, N_INITCLIENT, client.CN, client.Name, client.Team, client.PlayerModel)
	if client.GameState.State == CS_SPECTATOR {
		client.sendToAllOthers(true, 1, N_SPECTATOR, client.CN, 1)
	}
}

// Informs all other clients that a client left the game.
func (client *Client) informOthersOfDisconnect(reason DisconnectReason) {
	client.sendToAllOthers(true, 1, N_LEAVE, client.CN)
	// TOOD: send a server message with the disconnect reason in case it's not a normal leave
}

func (client *Client) spawn() {
	client.GameState.reset()
	client.GameState.spawn(state.GameMode)
	client.GameState.LifeSequence = (client.GameState.LifeSequence + 1) % 128

	parts := []interface{}{N_SPAWNSTATE, client.CN, client.GameState.LifeSequence, client.GameState.Health, client.GameState.MaxHealth, client.GameState.Armour, client.GameState.ArmourType, client.GameState.SelectedGun}
	for _, gn := range GunsWithAmmo {
		parts = append(parts, client.GameState.Ammo[gn])
	}
	client.send(true, 1, parts)

	client.GameState.LastSpawn = state.TimeLeft
}

// Resets the client object. Keeps the client's CN, so low CNs can be reused.
func (client *Client) reset() {
	log.Println("reset:", client.CN)

	client.Name = ""
	client.PlayerModel = -1
	client.Joined = false
	client.HasToAuthForConnect = false
	client.ReasonWhyAuthNeeded = DISC_NONE
	client.AI = false
	client.AISkill = -1
	client.InUse = false
	client.SessionId = rng.Int31()

	client.GameState.reset()
}
