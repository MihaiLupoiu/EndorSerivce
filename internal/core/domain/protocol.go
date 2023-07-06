package domain

import (
	"fmt"
	"sort"
	"strings"
)

type ProtocolType string

const (
	ClosestEnemies  ProtocolType = "closest-enemies"
	FurthestEnemies ProtocolType = "furthest-enemies"
	AssistAllies    ProtocolType = "assist-allies"
	AvoidCrossfire  ProtocolType = "avoid-crossfire"
	PrioritizeMech  ProtocolType = "prioritize-mech"
	AvoidMech       ProtocolType = "avoid-mech"
)

var (
	capabilitiesMap = map[string]ProtocolType{
		"closest-enemies":  ClosestEnemies,
		"furthest-enemies": FurthestEnemies,
		"assist-allies":    AssistAllies,
		"avoid-crossfire":  AvoidCrossfire,
		"prioritize-mech":  PrioritizeMech,
		"avoid-mech":       AvoidMech,
	}
)

// ParseStringToProtocolType parses a string representation of a protocol and
// returns the corresponding ProtocolType.
// If the string is not a valid protocol, an error is returned.
func ParseStringToProtocolType(str string) (ProtocolType, error) {
	c, ok := capabilitiesMap[strings.ToLower(str)]
	if !ok {
		return c, fmt.Errorf(`cannot parse:[%s] as ProtocolType`, str)
	}
	return c, nil
}

// Protocol is the interface that defines the methods for applying a protocol to a list of scans.
type Protocol interface {
	apply(scans []*Scan) []*Scan
}

// ProtocolDistanceLimit is a protocol that filters scans based on a maximum distance limit.
type ProtocolDistanceLimit struct {
	MaxDistance float64
}

// apply applies the ProtocolDistanceLimit to the provided scans and filters out scans beyond the maximum distance limit.
func (p ProtocolDistanceLimit) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	filteredScans := []*Scan{}

	for _, scan := range scans {
		enemyDistance := scan.Coordinates.GetDistance()
		if enemyDistance <= p.MaxDistance {
			filteredScans = append(filteredScans, scan)
		}
	}

	return filteredScans
}

// ProtocolClosestEnemies is a protocol that sorts scans based on the distance to the enemies in ascending order.
type ProtocolClosestEnemies struct{}

// apply applies the ProtocolClosestEnemies to the provided scans and sorts them based on the distance to enemies in ascending order.
func (p ProtocolClosestEnemies) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	sort.Slice(scans, func(i, j int) bool {
		distanceI := scans[i].Coordinates.GetDistance()
		distanceJ := scans[j].Coordinates.GetDistance()
		return distanceI < distanceJ
	})

	return scans
}

// ProtocolFurthestEnemies is a protocol that sorts scans based on the distance to the enemies in descending order.
type ProtocolFurthestEnemies struct{}

// apply applies the ProtocolFurthestEnemies to the provided scans and sorts them based on the distance to enemies in descending order.
func (p ProtocolFurthestEnemies) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	sort.Slice(scans, func(i, j int) bool {
		distanceI := scans[i].Coordinates.GetDistance()
		distanceJ := scans[j].Coordinates.GetDistance()
		return distanceI > distanceJ
	})

	return scans
}

// ProtocolAssistAllies is a protocol that filters scans to only include scans with allies present.
type ProtocolAssistAllies struct{}

// apply applies the ProtocolAssistAllies to the provided scans and filters out scans without any allies present.
func (p ProtocolAssistAllies) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	var safeScans []*Scan

	for _, scan := range scans {
		if scan.Allies > 0 {
			safeScans = append(safeScans, scan)
		}
	}

	return safeScans
}

// ProtocolAvoidCrossfire is a protocol that filters scans to only include scans without any allies present.
type ProtocolAvoidCrossfire struct{}

// apply applies the ProtocolAvoidCrossfire to the provided scans and filters out scans with any allies present.
func (p ProtocolAvoidCrossfire) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	var safeScans []*Scan

	for _, scan := range scans {
		if scan.Allies == 0 {
			safeScans = append(safeScans, scan)
		}
	}

	return safeScans
}

// ProtocolPrioritizeMech is a protocol that prioritizes scans with mech enemies and includes any other enemy type if no mech enemies are found.
type ProtocolPrioritizeMech struct{}

// apply applies the ProtocolPrioritizeMech to the provided scans and filters out scans that do not match the prioritization criteria.
func (p ProtocolPrioritizeMech) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	var prioritizedScans []*Scan

	for _, scan := range scans {
		if scan.Enemies != nil && scan.Enemies.Type == Mech {
			prioritizedScans = append(prioritizedScans, scan)
		}
	}

	if len(prioritizedScans) == 0 {
		// If no mech enemies found, include any other enemy type
		for _, scan := range scans {
			if scan.Enemies != nil && scan.Enemies.Type != Mech {
				prioritizedScans = append(prioritizedScans, scan)
			}
		}
	}

	return prioritizedScans
}

// ProtocolAvoidMech is a protocol that filters out scans with mech enemies.
type ProtocolAvoidMech struct{}

// apply applies the ProtocolAvoidMech to the provided scans and filters out scans with mech enemies.
func (p ProtocolAvoidMech) apply(scans []*Scan) []*Scan {
	if len(scans) == 0 {
		return scans
	}

	var safeScans []*Scan

	for _, scan := range scans {
		if scan.Enemies != nil && scan.Enemies.Type != Mech {
			safeScans = append(safeScans, scan)
		}
	}

	return safeScans
}

// GetProtocols returns a slice of Protocol instances based on the provided protocol types.
func GetProtocols(protocolTypes []ProtocolType) []Protocol {
	var protocols []Protocol

	// TODO: [Improvement] verify we don't have conflicting Protocol types.

	// TODO: [Improvement] make max distance configurable.
	// Add default distance limit of 100.
	protocols = append(protocols, ProtocolDistanceLimit{MaxDistance: 100})

	for _, protocolType := range protocolTypes {
		switch protocolType {
		case ClosestEnemies:
			protocols = append(protocols, ProtocolClosestEnemies{})
		case FurthestEnemies:
			protocols = append(protocols, ProtocolFurthestEnemies{})
		case AssistAllies:
			protocols = append(protocols, ProtocolAssistAllies{})
		case AvoidCrossfire:
			protocols = append(protocols, ProtocolAvoidCrossfire{})
		case PrioritizeMech:
			protocols = append(protocols, ProtocolPrioritizeMech{})
		case AvoidMech:
			protocols = append(protocols, ProtocolAvoidMech{})
		}
	}

	return protocols
}

// ApplyProtocols applies the specified protocols to the provided scans and returns the resulting scans.
func ApplyProtocols(scans []*Scan, protocol ...Protocol) []*Scan {
	result := make([]*Scan, len(scans))
	copy(result, scans)
	for _, proto := range protocol {
		result = proto.apply(result)
	}
	return result
}
