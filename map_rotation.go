package main

type MapRotation []string

// temporary set of maps used in development phase
var mapRotation = MapRotation{
	"hashi",
	"ot",
	"turbine",
	"shiva",
	"complex",
}

func (mr MapRotation) nextMap(mapName string) string {
	for i, m := range mr {
		if m == mapName {
			return mr[(i+1)%len(mr)]
		}
	}

	// current map wasn't found in map rotation, return random map in rotation
	return mr[rng.Intn(len(mr))]
}
