package domain

type Scan struct {
	Coordinates *Coordinate
	Enemies     *Enemy
	Allies      int
}

type Radar struct {
	Protocols []ProtocolType
	Scan      []*Scan
}
