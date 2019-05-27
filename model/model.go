package model

import "crypto/rand"

type Id []byte

func NewId() Id {
	ret := make(Id, 20)
	if _, err := rand.Read(ret); err != nil {
		panic(err)
	}
	return ret
}

// ssl grade model
type Grade struct {
	Score int
	Name  string
}

// ssl grades supported
var Grades = map[string]int{"A+": 1, "A-": 2, "A": 3, "B": 4, "C": 5, "D": 6, "E": 7, "F": 8, "T": 9, "M": 10}
