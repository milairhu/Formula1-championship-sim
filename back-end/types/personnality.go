package types

import (
	"math/rand"
)

type Personality struct {
	TraitsValue map[string]int `json:"personality"` // Dictionnary of traits
}

func NewPersonality(traitsValue map[string]int) *Personality {
	return &Personality{
		TraitsValue: traitsValue,
	}
}

func NewPersonalityRandom() *Personality {
	m := make(map[string]int)
	m["Aggressivity"] = rand.Intn(5) + 1
	m["Confidence"] = rand.Intn(5) + 1
	m["Docility"] = rand.Intn(5) + 1
	m["Concentration"] = rand.Intn(5) + 1
	return NewPersonality(m)
}
