package handler

import (
	"fmt"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/domain"
)

// This can be changed with an OpenAPI spec for example

type Coordinate struct {
	X *int `json:"x" validate:"required,min=0"`
	Y *int `json:"y" validate:"required,min=0"`
}

type Enemy struct {
	Type   string `json:"type" validate:"required"`
	Number *int   `json:"number" validate:"required,min=0"`
}

type Scan struct {
	Coordinates *Coordinate `json:"coordinates" validate:"required"`
	Enemies     *Enemy      `json:"enemies" validate:"required"`
	Allies      int         `json:"allies"`
}

type AttackRequest struct {
	Protocols []string `json:"protocols" validate:"required,dive,required"`
	Scan      []*Scan  `json:"scan" validate:"required,dive,required"`
}

// TODO: Review how to improve the convertion of the data
func (rq AttackRequest) ConvertToAttackDataModel() (*domain.Radar, error) {
	// Transforming data to the domain model.
	var attack = &domain.Radar{}

	for _, protocol := range rq.Protocols {
		proto, err := domain.ParseStringToProtocolType(protocol)
		if err != nil {
			return nil, fmt.Errorf("unable to process protocol type. err: %s", err)
		}
		attack.Protocols = append(attack.Protocols, proto)
	}

	for _, scan := range rq.Scan {
		coor := domain.NewCoordinates(*scan.Coordinates.X, *scan.Coordinates.Y)
		enemyType, err := domain.ParseStringToEnemyType(scan.Enemies.Type)
		if err != nil {
			return nil, fmt.Errorf("unable to process enemy type. err: %s", err)
		}
		enemy := &domain.Enemy{
			Type:   enemyType,
			Number: *scan.Enemies.Number,
		}

		var sc = &domain.Scan{
			Coordinates: coor,
			Enemies:     enemy,
			Allies:      scan.Allies,
		}
		attack.Scan = append(attack.Scan, sc)
	}
	return attack, nil
}

type AttackReportResponse struct {
	Casualties int         `json:"casualties"`
	Generation int         `json:"generation"`
	Target     *Coordinate `json:"target" validate:"required"`
}
