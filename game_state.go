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

// Sets GameState properties to the initial values depending on the mode.
func (gs *GameState) spawn(mode GameMode) {
	gs.QuadTimeLeft = 0
	gs.GunReloadTime = 0
	gs.State = CS_ALIVE

	switch mode {
	case GM_EFFIC, GM_EFFICTEAM, GM_EFFICCTF, GM_EFFICCOLLECT, GM_EFFICPROTECT, GM_EFFICHOLD:
		gs.Health = 100
		gs.MaxHealth = 100
		gs.Armour = 100
		gs.ArmourType = ARMOUR_GREEN
		gs.SelectedGun = GUN_MINIGUN
		gs.Ammo = SpawnAmmo[GM_EFFIC]

	case GM_INSTA, GM_INSTATEAM, GM_INSTACTF, GM_INSTACOLLECT, GM_INSTAPROTECT, GM_INSTAHOLD:
		gs.Health = 1
		gs.MaxHealth = 1
		gs.Armour = 0
		gs.ArmourType = ARMOUR_BLUE
		gs.SelectedGun = GUN_RIFLE
		gs.Ammo = SpawnAmmo[GM_INSTA]
	}
}

// Resets a player's game state.
func (gs *GameState) reset() {
	gs.Position = PlayerPosition{}
	gs.BufferedPackets = []Packet{}
	gs.HasReliablePacket = false

	if gs.State != CS_SPECTATOR {
		gs.State = CS_DEAD
	}
	gs.MaxHealth = 0
	gs.Tokens = 0

	gs.LastSpawn = 0
	gs.LifeSequence = 0
	gs.LastShot = 0
	gs.LastDeath = 0

	gs.Frags = 0
	gs.Deaths = 0
	gs.Teamkills = 0
	gs.ShotDamage = 0
	gs.Damage = 0
	gs.Flags = 0
}
