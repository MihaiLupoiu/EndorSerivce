package services

import (
	"fmt"
	"math"
	"sync"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/adapters"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/domain"
	"go.uber.org/zap"
)

// MAX_GORUTINES defines the maximum number of goroutines to execute concurrently.
var MAX_GORUTINES = 1000 // TODO: Have a configurable value
var log *zap.SugaredLogger

// semaphore is used to limit concurrency to MAX_GORUTINES.
var semaphore chan struct{}

// Simple semaphore to control the total numer of gorutines executed.
// For production we should have a hard limit on the number of gorutines we want in a way
// developers should not need to think much about it to  avoid loosing control over how
// many gorutines are executed in parralel.
// We can create a custom ParallelIonCannon and have that be in control. But in the real
// world scenario, we will query different endpoints with different parameters.

// EndorService represents the Endor service.
type EndorService struct {
	ionCannons []adapters.IonCannon
}

// NewEndorService creates a new instance of the EndorService.
func NewEndorService(logger *zap.SugaredLogger, ionCanons []adapters.IonCannon) *EndorService {
	log = logger
	semaphore = make(chan struct{}, MAX_GORUTINES)

	return &EndorService{
		ionCannons: ionCanons,
	}
}

// Attack performs the attack action on the specified target.
func (m *EndorService) Attack(attack *domain.Radar) (*domain.Report, error) {
	// We only have one action to make, so making a more complex structure does not make sense for now.
	listOfProtocols := domain.GetProtocols(attack.Protocols)
	targets := domain.ApplyProtocols(attack.Scan, listOfProtocols...)

	var finalTarget *domain.Coordinate
	var numberOfEmemies int
	if len(targets) > 0 {
		finalTarget = targets[0].Coordinates
		numberOfEmemies = targets[0].Enemies.Number
	} else {
		return nil, fmt.Errorf("not valid target encountered")
	}

	cas, gen, err := fire(finalTarget.X, finalTarget.Y, numberOfEmemies, m.ionCannons)
	if err != nil {
		return nil, err
	}

	report := &domain.Report{
		Target:     finalTarget,
		Casualties: cas,
		Generation: gen,
	}
	return report, nil
}

// fire fires the ion cannons at the specified target coordinates.
// It checks the status of all ion cannons concurrently and selects the one with the lowest generation to fire.
// If no ion cannons are available, it returns an error.
func fire(x, y, enemies int, ionCannons []adapters.IonCannon) (casualties int, generation int, err error) {
	type ChannelResponse struct {
		Available       bool
		Generation      int
		IonCannonClient adapters.IonCannon
		Err             error
	}

	responses := make(chan ChannelResponse, len(ionCannons))
	var wgCanon sync.WaitGroup

	// Query status of all ion cannons concurrently
	for _, c := range ionCannons {
		semaphore <- struct{}{}
		wgCanon.Add(1)
		go func(c adapters.IonCannon) {
			defer func() {
				defer wgCanon.Done()
				<-semaphore
			}()
			res, err := c.CheckStatus()
			if err != nil {
				log.Errorf("Failed to check status: %v\n", err)
				responses <- ChannelResponse{Available: false, Generation: 0, Err: err}
				return
			}

			responses <- ChannelResponse{
				Available:       res.Available,
				Generation:      res.Generation,
				IonCannonClient: c,
				Err:             nil,
			}
		}(c)
	}

	wgCanon.Wait()
	close(responses)

	// Get the lowest generation of ion cannon available to fire
	lowestGeneration := math.MaxInt
	var lowestGenerationIonCannon adapters.IonCannon

	for response := range responses {
		if response.Err != nil {
			continue
		}

		if response.Available {
			if response.Generation < lowestGeneration {
				lowestGeneration = response.Generation
				lowestGenerationIonCannon = response.IonCannonClient
			}
		}
	}

	if lowestGenerationIonCannon == nil {
		return 0, 0, fmt.Errorf("failed to fire. No available ion cannons")
	}

	casualties, generation, err = lowestGenerationIonCannon.FireCommand(x, y, enemies)
	if err != nil {
		log.Errorf("Failed to fire command: %v\n", err)
		return 0, 0, err
	}

	return casualties, generation, nil

}
