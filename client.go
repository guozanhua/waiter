package main

import (
	"./enet"
	"log"
)

// player's cn
type ClientNumber int32

// describes a client's level of privilege
type ClientPrivilege int32

const (
	PRIV_NONE ClientPrivilege = iota
	PRIV_MASTER
	PRIV_AUTH
	PRIV_ADMIN
)

// Describes a client
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
	Peer                enet.Peer
	SessionId           int32
}

// Connects an ENet peer to a client object. If no unused client object can be found, a new one is created and added to the global set of clients
func addClient(peer enet.Peer) *Client {
	var client *Client

	// re-use unused client object with low cn
	for _, client = range clients {
		if !client.InUse {
			client.InUse = true
			return client
		}
	}

	client = &Client{
		CN: ClientNumber(len(clients)),
	}

	client.InUse = true
	client.Peer = peer
	client.SessionId = rng.Int31()
	client.Team = "good" // TODO: select weaker team

	client.GameState = NewGameState(5)

	// first client connected, start new game
	if len(clients) == 0 {
		go countDown(TEN_MINUTES, interruptTickerChannel, intermissionChannel)
	}

	clients[client.CN] = client

	log.Println("added client:", client, "with peer", client.Peer)

	return client
}

// send a packet to a client (reliable if desired) over the specified channel
func (client *Client) send(p *Packet, reliable bool, channel uint8) {
	log.Println(p.buf, "→", client.Peer.Address.String())

	var flags enet.PacketFlag
	if reliable {
		flags |= enet.PACKET_FLAG_RELIABLE
	}

	client.Peer.Send(p.buf, flags, channel)
}

// sends basic server info to the client
func (client *Client) sendServerInfo() {
	sendf(client, true, 1, N_SERVINFO, client.CN, PROTOCOL_VERSION, client.SessionId, 0, config.ServerDescription, config.ServerAuthDomains[0])
}

// sends 'welcome' information to a newly joined client, like map, mode, time left, other players, etc.
func (client *Client) sendWelcome() {
	parts := []interface{}{N_WELCOME}

	// send currently played map
	parts = append(parts, N_MAPCHANGE, state.Map, state.GameMode, state.NotGotItems)

	// send time left in this round
	parts = append(parts, N_TIMELEFT, state.TimeLeft/1000)

	// send list of clients which have privilege other than PRIV_NONE and their respecitve privilege level
	//parts = append(parts, N_CURRENTMASTER, state.MasterMode)

	// tell the client what team he was put in by the server
	parts = append(parts, N_SETTEAM, client.CN, client.Team, int32(-1))

	// tell the client how to spawn (what health, what armour, what weapons, what ammo, etc.)
	// TODO: handle spectators (locked mode) → spawn as dead (?)
	parts = append(parts, N_SPAWNSTATE, client.CN, client.GameState.LifeSequence, client.GameState.Health, client.GameState.MaxHealth, client.GameState.Armour, client.GameState.ArmourType, client.GameState.SelectedGun)
	for _, gn := range GunsWithAmmo {
		parts = append(parts, client.GameState.Ammo[gn])
	}

	// send other players' state (frags, flags, etc.)
	parts = append(parts, N_RESUME)
	for _, c := range clients {
		log.Println(c)
		if c == client || !client.InUse {
			continue
		}

		parts = append(parts, c.CN, c.GameState.State, c.GameState.Frags, c.GameState.Flags, c.GameState.QuadTimeLeft, c.GameState.LifeSequence, c.GameState.Health, c.GameState.MaxHealth, c.GameState.Armour, c.GameState.ArmourType, c.GameState.SelectedGun)
		for _, gn := range GunsWithAmmo {
			parts = append(parts, c.GameState.Ammo[gn])
		}
	}
	parts = append(parts, int32(-1))

	// send other client's state (name, team, playermodel)
	for _, c := range clients {
		if c == client || !client.InUse {
			continue
		}
		parts = append(parts, N_INITCLIENT, c.CN, c.Name, c.Team, c.PlayerModel)
	}

	log.Println(parts)

	sendf(client, true, 1, parts...)
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

	log.Printf("join: %s (%d)\n", name, client.CN)

	// send welcome packet
	client.sendWelcome()

	// inform other clients that there is a new client
	for _, c := range clients {
		if c == client {
			continue
		}
		c.informOfNewClient(client)
	}

	return
}

// Informs the client that another client joined the game.
func (client *Client) informOfNewClient(otherClient *Client) {
	sendf(client, true, 1, N_INITCLIENT, otherClient.CN, otherClient.Name, otherClient.Team, otherClient.PlayerModel)
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
	for _, c := range clients {
		if c == client {
			continue
		}
		c.informOfDisconnectingClient(client, reason)
	}

	if reason != DISC_NONE {
		log.Printf("disconnected: %s (%d) - %s", DisconnectReasons[reason])
	}

	client.Peer.Disconnect(uint32(reason))

	client.reset()
}

// Informs the client that a client left the game
func (client *Client) informOfDisconnectingClient(otherClient *Client, reason DisconnectReason) {
	sendf(client, true, 1, N_LEAVE, otherClient.CN)

	// TOOD: send a server message with the disconnect reason in case it's not a normal leave
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
	client.SessionId = -1

	client.GameState.reset()
}
