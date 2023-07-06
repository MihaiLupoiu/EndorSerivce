package services

import (
	"testing"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/adapters"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/common/logger"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/domain"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestEndorService_Attack(t *testing.T) {
	// TODO: create abstraction of logger
	log := logger.NewLogger(logger.DEBUG, false)

	// Create an instance of the mock
	mockIonCannonV1 := &mocks.IonCannonClientMock{
		// Mock the CheckStatus function as needed
		CheckStatusFunc: func() (*domain.IonCannon, error) {
			return &domain.IonCannon{Available: true, Generation: 1}, nil
		},

		// Mock the FireCommand function as needed
		FireCommandFunc: func(targetX int, targetY int, enemies int) (casualties int, generation int, err error) {
			return enemies, 1, nil
		},
	}
	mockIonCannonV2 := &mocks.IonCannonClientMock{
		// Mock the CheckStatus function as needed
		CheckStatusFunc: func() (*domain.IonCannon, error) {
			return &domain.IonCannon{Available: true, Generation: 2}, nil
		},

		// Mock the FireCommand function as needed
		FireCommandFunc: func(targetX int, targetY int, enemies int) (casualties int, generation int, err error) {
			return enemies, 2, nil
		},
	}

	// Create an instance of the EndorService with the mock IonCannon
	endorService := NewEndorService(log, []adapters.IonCannon{mockIonCannonV1, mockIonCannonV2})

	// Create a sample attack
	attack := &domain.Radar{
		Protocols: []domain.ProtocolType{domain.ClosestEnemies},
		Scan: []*domain.Scan{
			{
				Coordinates: domain.NewCoordinates(10, 20),
				Enemies: &domain.Enemy{
					Type:   domain.Soldier,
					Number: 5,
				},
				Allies: 0,
			},
		},
	}

	// Call the Attack function
	report, err := endorService.Attack(attack)

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 5, report.Casualties)
	assert.Equal(t, 1, report.Generation)

	// Assert the function calls on the mock
	assert.Len(t, mockIonCannonV1.CheckStatusCallData, 1)
	assert.Len(t, mockIonCannonV1.FireCommandCallData, 1)
	assert.Len(t, mockIonCannonV2.CheckStatusCallData, 1)
	assert.Len(t, mockIonCannonV2.FireCommandCallData, 0)
	assert.Equal(t, 10, mockIonCannonV1.FireCommandCallData[0].TargetX)
	assert.Equal(t, 20, mockIonCannonV1.FireCommandCallData[0].TargetY)
	assert.Equal(t, 5, mockIonCannonV1.FireCommandCallData[0].Enemies)
}
