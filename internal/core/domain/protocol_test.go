package domain

import (
	"reflect"
	"testing"
)

func TestGetProtocols(t *testing.T) {
	testCases := []struct {
		protocolTypes []ProtocolType
		expected      []Protocol
	}{
		{
			protocolTypes: []ProtocolType{ClosestEnemies},
			expected: []Protocol{
				ProtocolDistanceLimit{MaxDistance: 100},
				ProtocolClosestEnemies{},
			},
		},
		{
			protocolTypes: []ProtocolType{ClosestEnemies, AvoidCrossfire},
			expected: []Protocol{
				ProtocolDistanceLimit{MaxDistance: 100},
				ProtocolClosestEnemies{},
				ProtocolAvoidCrossfire{},
			},
		},
		{
			protocolTypes: []ProtocolType{FurthestEnemies, PrioritizeMech},
			expected: []Protocol{
				ProtocolDistanceLimit{MaxDistance: 100},
				ProtocolFurthestEnemies{},
				ProtocolPrioritizeMech{},
			},
		},
		// TODO: Add more test cases as needed
	}

	for _, testCase := range testCases {
		result := GetProtocols(testCase.protocolTypes)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Unexpected result. Expected: %v, Got: %v", testCase.expected, result)
		}
	}
}

func TestApplyProtocols(t *testing.T) {
	// Define sample scans for testing
	scans := []*Scan{
		{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
		{Coordinates: NewCoordinates(2, 2), Enemies: &Enemy{Type: Mech, Number: 1}, Allies: 0},
		{Coordinates: NewCoordinates(4, 4), Enemies: &Enemy{Type: Mech, Number: 2}, Allies: 1},
		{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
		{Coordinates: NewCoordinates(105, 105), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 1},
	}

	// Define test cases with different combinations of protocols
	testCases := []struct {
		name      string
		protocols []Protocol
		expected  []*Scan
	}{
		{
			name:      "test distance filter",
			protocols: []Protocol{ProtocolDistanceLimit{100}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
				{Coordinates: NewCoordinates(2, 2), Enemies: &Enemy{Type: Mech, Number: 1}, Allies: 0},
				{Coordinates: NewCoordinates(4, 4), Enemies: &Enemy{Type: Mech, Number: 2}, Allies: 1},
				{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
			},
		},
		{
			name:      "test ProtocolClosestEnemies",
			protocols: []Protocol{ProtocolClosestEnemies{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
				{Coordinates: NewCoordinates(2, 2), Enemies: &Enemy{Type: Mech, Number: 1}, Allies: 0},
				{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
				{Coordinates: NewCoordinates(4, 4), Enemies: &Enemy{Type: Mech, Number: 2}, Allies: 1},
				{Coordinates: NewCoordinates(105, 105), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 1},
			},
		},
		{
			name:      "test ProtocolFurthestEnemies",
			protocols: []Protocol{ProtocolFurthestEnemies{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(105, 105), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 1},
				{Coordinates: NewCoordinates(4, 4), Enemies: &Enemy{Type: Mech, Number: 2}, Allies: 1},
				{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
				{Coordinates: NewCoordinates(2, 2), Enemies: &Enemy{Type: Mech, Number: 1}, Allies: 0},
				{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
			},
		},
		{
			name:      "test ProtocolAssistAllies",
			protocols: []Protocol{ProtocolAssistAllies{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(4, 4), Enemies: &Enemy{Type: Mech, Number: 2}, Allies: 1},
				{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
				{Coordinates: NewCoordinates(105, 105), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 1},
			},
		},
		{
			name:      "test ProtocolAvoidCrossfire",
			protocols: []Protocol{ProtocolAvoidCrossfire{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
				{Coordinates: NewCoordinates(2, 2), Enemies: &Enemy{Type: Mech, Number: 1}, Allies: 0},
			},
		},
		{
			name:      "test ProtocolPrioritizeMech",
			protocols: []Protocol{ProtocolPrioritizeMech{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(2, 2), Enemies: &Enemy{Type: Mech, Number: 1}, Allies: 0},
				{Coordinates: NewCoordinates(4, 4), Enemies: &Enemy{Type: Mech, Number: 2}, Allies: 1},
			},
		},
		{
			name:      "test ProtocolAvoidMech",
			protocols: []Protocol{ProtocolAvoidMech{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
				{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
				{Coordinates: NewCoordinates(105, 105), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 1},
			},
		},
		{
			name:      "test ProtocolClosestEnemies and ProtocolAvoidMech",
			protocols: []Protocol{ProtocolClosestEnemies{}, ProtocolAvoidMech{}},
			expected: []*Scan{
				{Coordinates: NewCoordinates(1, 1), Enemies: &Enemy{Type: Soldier, Number: 1}, Allies: 2},
				{Coordinates: NewCoordinates(3, 3), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 0},
				{Coordinates: NewCoordinates(105, 105), Enemies: &Enemy{Type: Soldier, Number: 2}, Allies: 1},
			},
		},
		// Add more test cases as needed
	}

	for _, testCase := range testCases {
		result := ApplyProtocols(scans, testCase.protocols...)

		if len(result) != len(testCase.expected) {
			t.Errorf("Unexpected result size in test %s. Expected: %v, Got: %v", testCase.name, testCase.expected, result)
			continue
		}
		// Compare the result with the expected value
		for i := 0; i < len(result); i++ {
			if result[i].Coordinates.X != testCase.expected[i].Coordinates.X {
				t.Errorf("Unexpected X coordinate at index %d. Expected: %d, Got: %d", i, testCase.expected[i].Coordinates.X, result[i].Coordinates.X)
			}
			if result[i].Coordinates.Y != testCase.expected[i].Coordinates.Y {
				t.Errorf("Unexpected Y coordinate at index %d. Expected: %d, Got: %d", i, testCase.expected[i].Coordinates.Y, result[i].Coordinates.Y)
			}
			if result[i].Enemies.Type != testCase.expected[i].Enemies.Type {
				t.Errorf("Unexpected enemy type at index %d. Expected: %s, Got: %s", i, testCase.expected[i].Enemies.Type, result[i].Enemies.Type)
			}
			if result[i].Enemies.Number != testCase.expected[i].Enemies.Number {
				t.Errorf("Unexpected enemy number at index %d. Expected: %d, Got: %d", i, testCase.expected[i].Enemies.Number, result[i].Enemies.Number)
			}
			if result[i].Allies != testCase.expected[i].Allies {
				t.Errorf("Unexpected number of allies at index %d. Expected: %d, Got: %d", i, testCase.expected[i].Allies, result[i].Allies)
			}
		}
	}
}

func TestProtocolAvoidMech_apply(t *testing.T) {
	type args struct {
		scans []*Scan
	}
	tests := []struct {
		name string
		p    ProtocolAvoidMech
		args args
		want []*Scan
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProtocolAvoidMech{}
			if got := p.apply(tt.args.scans); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtocolAvoidMech.apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
