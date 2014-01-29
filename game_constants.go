package main

// armour

type ArmourType int32

const (
	ARMOUR_BLUE ArmourType = iota
	ARMOUR_GREEN
	ARMOUR_YELLOW
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
	WPN_SAW:             Weapon{S_SAW, 250, 50, 0, 0, 0, 14, 1, 80, 0, 0},
	WPN_SHOTGUN:         Weapon{S_SHOTGUN, 1400, 10, 400, 0, 20, 1024, 20, 80, 0, 0},
	WPN_MINIGUN:         Weapon{S_MINIGUN, 100, 30, 100, 0, 7, 1024, 1, 80, 0, 0},
	WPN_ROCKETLAUNCHER:  Weapon{S_ROCKETLAUNCH, 800, 120, 0, 320, 10, 1024, 1, 160, 40, 0},
	WPN_RIFLE:           Weapon{S_RIFLE, 1500, 100, 0, 0, 30, 2048, 1, 80, 0, 0},
	WPN_GRENADELAUNCHER: Weapon{S_GRENADELAUNCH, 600, 90, 0, 200, 10, 1024, 1, 250, 45, 1500},
	WPN_PISTOL:          Weapon{S_PISTOL, 500, 35, 50, 0, 7, 1024, 1, 80, 0, 0},
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
}

// entities (especially pickups)

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
	PICKUP_SHOTGUN // 8
	PICKUP_MINIGUN
	PICKUP_ROCKETLAUNCHER
	PICKUP_RIFLE
	PICKUP_GRENADELAUNCHER
	PICKUP_PISTOL
	PICKUP_HEALTH
	PICKUP_BOOST
	PICKUP_GREENARMOUR
	PICKUP_YELLOWARMOUR
	PICKUP_QUAD
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
	PICKUP_SHOTGUN:         PickUp{S_PICKUP_AMMO, 10},
	PICKUP_MINIGUN:         PickUp{S_PICKUP_AMMO, 20},
	PICKUP_ROCKETLAUNCHER:  PickUp{S_PICKUP_AMMO, 5},
	PICKUP_RIFLE:           PickUp{S_PICKUP_AMMO, 5},
	PICKUP_GRENADELAUNCHER: PickUp{S_PICKUP_AMMO, 10},
	PICKUP_PISTOL:          PickUp{S_PICKUP_AMMO, 30},
	PICKUP_HEALTH:          PickUp{S_PICKUP_HEALTH, 25},
	PICKUP_BOOST:           PickUp{S_PICKUP_HEALTH, 10},
	PICKUP_GREENARMOUR:     PickUp{S_PICKUP_ARMOUR, 100},
	PICKUP_YELLOWARMOUR:    PickUp{S_PICKUP_ARMOUR, 200},
	PICKUP_QUAD:            PickUp{S_PICKUP_QUAD, 20000},
}

// sounds

// Used to tell clients what sound to play
type SoundNumber int32

const (
	S_JUMP SoundNumber = iota
	S_LAND
	S_RIFLE
	S_SAW
	S_SHOTGUN
	S_MINIGUN
	S_ROCKETLAUNCH
	S_RLHIT
	S_WEAPLOAD
	S_PICKUP_AMMO
	S_PICKUP_HEALTH
	S_PICKUP_ARMOUR
	S_PICKUP_QUAD
	S_ITEMSPAWN
	S_TELEPORT
	S_NOAMMO
	S_PUPOUT
	S_PAIN1
	S_PAIN2
	S_PAIN3
	S_PAIN4
	S_PAIN5
	S_PAIN6
	S_DIE1
	S_DIE2
	S_GRENADELAUNCH
	S_FEXPLODE
	S_SPLASH1
	S_SPLASH2
	S_GRUNT1
	S_GRUNT2
	S_RUMBLE
	S_PAINO
	S_PAINR
	S_DEATHR
	S_PAINE
	S_DEATHE
	S_PAINS
	S_DEATHS
	S_PAINB
	S_DEATHB
	S_PAINP
	S_PIGGR2
	S_PAINH
	S_DEATHH
	S_PAIND
	S_DEATHD
	S_PIGR1
	S_ICEBALL
	S_SLIMEBALL
	S_JUMPPAD
	S_PISTOL
	S_V_BASECAP
	S_V_BASELOST
	S_V_FIGHT
	S_V_BOOST
	S_V_BOOST10
	S_V_QUAD
	S_V_QUAD10
	S_V_RESPAWNPOINT
	S_FLAGPICKUP
	S_FLAGDROP
	S_FLAGRETURN
	S_FLAGSCORE
	S_FLAGRESET
	S_BURN
	S_CHAINSAW_ATTACK
	S_CHAINSAW_IDLE
	S_HIT
	S_FLAGFAIL
)

// master modes

type MasterMode int32

const (
	MM_AUTH MasterMode = -1
	MM_OPEN            = iota
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
