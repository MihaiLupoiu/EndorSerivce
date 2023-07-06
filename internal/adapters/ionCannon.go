package adapters

import "github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/domain"

// IonCannon Interface for service to use
type IonCannon interface {
	CheckStatus() (*domain.IonCannon, error)
	FireCommand(targetX int, targetY int, enemies int) (casualties int, generation int, err error)
}
