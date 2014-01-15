package main

// a player's position in the map/world
type PlayerPosition Packet

// client's state
type ClientState uint

const (
	CS_ALIVE ClientState = iota
	CS_DEAD
	CS_SPAWNING
	CS_LAGGED
	CS_EDITING
	CS_SPECTATOR
)

// The game state of a player
type GameState struct {
	// position of player
	Position PlayerPosition

	// buffered packets the client send which need to be sent to the other clients
	BufferedPackets   []Packet
	HasReliablePacket bool // wether one of the packets is important and nees to be sent reliably

	// fields that change at spawn
	State         ClientState
	Health        int32
	MaxHealth     int32
	Armour        int32
	ArmourType    ArmourType
	QuadTimeLeft  int32 // in milliseconds
	SelectedGun   GunNumber
	GunReloadTime int32
	Ammo          map[GunNumber]int32
	Tokens        int32 // skulls

	LastSpawn    int32
	LifeSequence int32
	LastShot     int32
	LastDeath    int32

	// fields that change at intermission
	Frags      int32
	Deaths     int32
	Teamkills  int32
	ShotDamage int32
	Damage     int32
	Flags      int32
}

// Returns a fresh game state depending on the game mode
func NewGameState(mode GameMode) GameState {
	gs := GameState{}

	switch mode {
	case GM_EFFIC, GM_EFFICTEAM:
		gs.Health = 100
		gs.MaxHealth = 100
		gs.Armour = 100
		gs.ArmourType = ARMOUR_GREEN
		gs.SelectedGun = GUN_MINIGUN

		gs.Ammo = map[GunNumber]int32{}
		baseAmmo(gs.Ammo)
		gs.Ammo[GUN_MINIGUN] /= 2
		gs.Ammo[GUN_PISTOL] = 0

	}

	return gs
}

// Resets a player's game state.
func (gs *GameState) reset() {
	if gs.State != CS_SPECTATOR {
		gs.State = CS_DEAD
	}
	gs.MaxHealth = 0
	gs.Tokens = 0

	gs.LastDeath = 0

	gs.Frags = 0
	gs.Deaths = 0
	gs.Teamkills = 0
	gs.ShotDamage = 0
	gs.Damage = 0
}
