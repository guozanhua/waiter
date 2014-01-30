package main

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
	Position Packet

	// fields that change at spawn
	State          ClientState
	Health         int32
	MaxHealth      int32
	Armour         int32
	ArmourType     ArmourType
	QuadTimeLeft   int32 // in milliseconds
	SelectedWeapon WeaponNumber
	GunReloadTime  int32
	Ammo           map[WeaponNumber]int32
	Tokens         int32 // skulls

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
	gs.Tokens = 0

	switch mode {
	case GM_EFFIC, GM_EFFICTEAM, GM_EFFICCTF, GM_EFFICCOLLECT, GM_EFFICPROTECT, GM_EFFICHOLD:
		gs.Health = 100
		gs.MaxHealth = 100
		gs.Armour = 100
		gs.ArmourType = ARMOUR_GREEN
		gs.SelectedWeapon = WPN_MINIGUN
		gs.Ammo = SpawnAmmo[GM_EFFIC]

	case GM_INSTA, GM_INSTATEAM, GM_INSTACTF, GM_INSTACOLLECT, GM_INSTAPROTECT, GM_INSTAHOLD:
		gs.Health = 1
		gs.MaxHealth = 1
		gs.Armour = 0
		gs.ArmourType = ARMOUR_BLUE
		gs.SelectedWeapon = WPN_RIFLE
		gs.Ammo = SpawnAmmo[GM_INSTA]
	}
}

func (gs *GameState) selectWeapon(selectedWeapon WeaponNumber) {
	if gs.State != CS_ALIVE {
		return
	}

	if selectedWeapon >= WPN_SAW && selectedWeapon <= WPN_PISTOL {
		gs.SelectedWeapon = selectedWeapon
	} else {
		gs.SelectedWeapon = WPN_PISTOL
	}
}

// Resets a player's game state.
func (gs *GameState) reset() {
	gs.Position = Packet{}

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
