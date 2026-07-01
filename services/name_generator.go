package services

import (
	"crypto/rand"
	"math/big"
	"strings"

	_ "embed"
)

//go:embed adjectives.txt
var adjectivesRaw string
var adjectives = strings.Fields(adjectivesRaw)

//go:embed nouns.txt
var nounsRaw string
var animals = strings.Fields(nounsRaw)

func GenerateTunnelName() (string, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	if err != nil {
		return "", err
	}

	j, err := rand.Int(rand.Reader, big.NewInt(int64(len(animals))))
	if err != nil {
		return "", err
	}

	return strings.ToLower(adjectives[i.Int64()]) + "-" + strings.ToLower(animals[j.Int64()]), nil
}
