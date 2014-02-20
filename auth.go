package main

import (
	"math/rand"
	"time"
)

var (
	rng *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	authRequestNumber uint32 = 1
)

const (
	MASTER_AUTH_REQUEST = "reqauth"
	MASTER_AUTH_CONFIRM = "confauth"
)

func tryAuth(name, domain string) bool {
	return false
}
