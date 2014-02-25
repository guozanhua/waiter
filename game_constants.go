package main

// armour

type ArmourType int32

const (
	ARM_BLUE ArmourType = iota
	ARM_GREEN
	ARM_YELLOW
)

// guns

type WeaponNumber int32

const (
	WPN_SAW WeaponNumber = iota
	WPN_SHOTGUN
	WPN_MINIGUN
	WPN_ROCKETLAUNCHER
	WPN_RIFLE
	WPN_GRENADELAUNCHER
	WPN_PISTOL
)

var WeaponsWithAmmo []WeaponNumber = []WeaponNumber{WPN_SHOTGUN, WPN_MINIGUN, WPN_ROCKETLAUNCHER, WPN_RIFLE, WPN_GRENADELAUNCHER, WPN_PISTOL}

type Weapon struct {
	Sound           SoundNumber
	ReloadTime      int32
	Damage          int32
	Spread          int32
	ProjectileSpeed int32
	Recoil          int32
	Range           int32
	Rays            int32
	HitPush         int32
	ExplosionRadius int32
	TimeToLive      int32
}

var Weapons map[WeaponNumber]Weapon = map[WeaponNumber]Weapon{
	WPN_SAW:             Weapon{SND_SAW, 250, 50, 0, 0, 0, 14, 1, 80, 0, 0},
	WPN_SHOTGUN:         Weapon{SND_SHOTGUN, 1400, 10, 400, 0, 20, 1024, 20, 80, 0, 0},
	WPN_MINIGUN:         Weapon{SND_MINIGUN, 100, 30, 100, 0, 7, 1024, 1, 80, 0, 0},
	WPN_ROCKETLAUNCHER:  Weapon{SND_ROCKETLAUNCH, 800, 120, 0, 320, 10, 1024, 1, 160, 40, 0},
	WPN_RIFLE:           Weapon{SND_RIFLE, 1500, 100, 0, 0, 30, 2048, 1, 80, 0, 0},
	WPN_GRENADELAUNCHER: Weapon{SND_GRENADELAUNCH, 600, 90, 0, 200, 10, 1024, 1, 250, 45, 1500},
	WPN_PISTOL:          Weapon{SND_PISTOL, 500, 35, 50, 0, 7, 1024, 1, 80, 0, 0},
}

// ammo sets

var SpawnAmmo map[GameMode]map[WeaponNumber]int32 = map[GameMode]map[WeaponNumber]int32{
	GM_EFFIC: map[WeaponNumber]int32{
		WPN_SHOTGUN:         20,
		WPN_MINIGUN:         20,
		WPN_ROCKETLAUNCHER:  10,
		WPN_RIFLE:           10,
		WPN_GRENADELAUNCHER: 20,
		WPN_PISTOL:          0,
	},
	GM_INSTA: map[WeaponNumber]int32{
		WPN_SHOTGUN:         0,
		WPN_MINIGUN:         0,
		WPN_ROCKETLAUNCHER:  0,
		WPN_RIFLE:           100,
		WPN_GRENADELAUNCHER: 0,
		WPN_PISTOL:          0,
	},
	// TODO: add all modes' spawn ammo sets
}

// entities (flags, bases, jumppads, pickups, teleports, teledests, etc.)

type EntityNumber int32

const (
	_ EntityNumber = iota
	LIGHT
	MAPMODEL
	PLAYERSTART
	ENVMAP
	PARTICLES
	MAPSOUND
	SPOTLIGHT
	PU_SHOTGUN // 8
	PU_MINIGUN
	PU_ROCKETLAUNCHER
	PU_RIFLE
	PU_GRENADELAUNCHER
	PU_PISTOL
	PU_HEALTH
	PU_BOOST
	PU_GREENARMOUR
	PU_YELLOWARMOUR
	PU_QUAD
	TELEPORT
	TELEDEST
	MONSTER
	CARROT
	JUMPPAD
	BASE
	RESPAWNPOINT
	BOX
	BARREL
	PLATFORM
	ELEVATOR
	FLAG
	MAXENTTYPES
)

type PickUp struct {
	Sound  SoundNumber
	Amount int32
}

var PickUps map[EntityNumber]PickUp = map[EntityNumber]PickUp{
	PU_SHOTGUN:         PickUp{SND_PU_AMMO, 10},
	PU_MINIGUN:         PickUp{SND_PU_AMMO, 20},
	PU_ROCKETLAUNCHER:  PickUp{SND_PU_AMMO, 5},
	PU_RIFLE:           PickUp{SND_PU_AMMO, 5},
	PU_GRENADELAUNCHER: PickUp{SND_PU_AMMO, 10},
	PU_PISTOL:          PickUp{SND_PU_AMMO, 30},
	PU_HEALTH:          PickUp{SND_PU_HEALTH, 25},
	PU_BOOST:           PickUp{SND_PU_HEALTH, 10},
	PU_GREENARMOUR:     PickUp{SND_PU_ARMOUR, 100},
	PU_YELLOWARMOUR:    PickUp{SND_PU_ARMOUR, 200},
	PU_QUAD:            PickUp{SND_PU_QUAD, 20000},
}

// sounds

// Sound numbers are used to tell clients what sound to play.
type SoundNumber int32

const (
	SND_JUMP SoundNumber = iota
	SND_LAND
	SND_RIFLE
	SND_SAW
	SND_SHOTGUN
	SND_MINIGUN
	SND_ROCKETLAUNCH
	SND_RLHIT
	SND_WEAPLOAD
	SND_PU_AMMO
	SND_PU_HEALTH
	SND_PU_ARMOUR
	SND_PU_QUAD
	SND_ITEMSPAWN
	SND_TELEPORT
	SND_NOAMMO
	SND_PUPOUT
	SND_PAIN1
	SND_PAIN2
	SND_PAIN3
	SND_PAIN4
	SND_PAIN5
	SND_PAIN6
	SND_DIE1
	SND_DIE2
	SND_GRENADELAUNCH
	SND_FEXPLODE
	SND_SPLASH1
	SND_SPLASH2
	SND_GRUNT1
	SND_GRUNT2
	SND_RUMBLE
	SND_PAINO
	SND_PAINR
	SND_DEATHR
	SND_PAINE
	SND_DEATHE
	SND_PAINS
	SND_DEATHS
	SND_PAINB
	SND_DEATHB
	SND_PAINP
	SND_PIGGR2
	SND_PAINH
	SND_DEATHH
	SND_PAIND
	SND_DEATHD
	SND_PIGR1
	SND_ICEBALL
	SND_SLIMEBALL
	SND_JUMPPAD
	SND_PISTOL
	SND_V_BASECAP
	SND_V_BASELOST
	SND_V_FIGHT
	SND_V_BOOST
	SND_V_BOOST10
	SND_V_QUAD
	SND_V_QUAD10
	SND_V_RESPAWNPOINT
	SND_FLAGPICKUP
	SND_FLAGDROP
	SND_FLAGRETURN
	SND_FLAGSCORE
	SND_FLAGRESET
	SND_BURN
	SND_CHAINSAW_ATTACK
	SND_CHAINSAW_IDLE
	SND_HIT
	SND_FLAGFAIL
)

// master modes

type MasterMode int32

const (
	MM_AUTH MasterMode = iota - 1
	MM_OPEN
	MM_VETO
	MM_LOCKED
	MM_PRIVATE
)

// game modes

type GameMode int32

const (
	GM_FFA GameMode = iota
	GM_COOPEDIT
	GM_TEAMPLAY
	GM_INSTA
	GM_INSTATEAM
	GM_EFFIC
	GM_EFFICTEAM
	GM_TACTICS
	GM_TACTICSTEAM
	GM_CAPTURE
	GM_REGENCAPTURE
	GM_CTF
	GM_INSTACTF
	GM_PROTECT
	GM_INSTAPROTECT
	GM_HOLD
	GM_INSTAHOLD
	GM_EFFICCTF
	GM_EFFICPROTECT
	GM_EFFICHOLD
	GM_COLLECT
	GM_INSTACOLLECT
	GM_EFFICCOLLECT
)

// disconnect reasons

type DisconnectReason uint32

const (
	DISC_NONE DisconnectReason = iota
	DISC_EOP
	DISC_LOCAL
	DISC_KICK
	DISC_MSGERR
	DISC_IPBAN
	DISC_PRIVATE
	DISC_MAXCLIENTS
	DISC_TIMEOUT
	DISC_OVERFLOW
	DISC_PASSWORD
	DISCNUM
)

var DisconnectReasons []string = []string{
	"",
	"end of packet",
	"server is in local mode",
	"kicked/banned",
	"message error",
	"ip is banned",
	"server is in private mode",
	"server full",
	"connection timed out",
	"overflow",
	"invalid password",
}
