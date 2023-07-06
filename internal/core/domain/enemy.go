package domain

import (
	"fmt"
	"strings"
)

// EnemyType represents the type of enemy.
type EnemyType string

const (
	Soldier EnemyType = "soldier"
	Mech    EnemyType = "mech"
)

type Enemy struct {
	Type   EnemyType
	Number int
}

var (
	EnemyMap = map[string]EnemyType{
		"soldier": Soldier,
		"mech":    Mech,
	}
)

// ParseStringToEnemyType parses a string and returns the corresponding EnemyType.
// It returns an error if the string cannot be parsed.
func ParseStringToEnemyType(str string) (EnemyType, error) {
	c, ok := EnemyMap[strings.ToLower(str)]
	if !ok {
		return c, fmt.Errorf(`cannot parse:[%s] as EnemyType`, str)
	}
	return c, nil
}
