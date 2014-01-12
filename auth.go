package main

import (
	"math/rand"
	"time"
)

var rng *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
