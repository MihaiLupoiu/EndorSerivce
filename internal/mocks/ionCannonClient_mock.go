package mocks

import (
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/domain"
)

type IonCannonClientMock struct {
	CheckStatusFunc     func() (*domain.IonCannon, error)
	CheckStatusCallData []struct{}
	FireCommandFunc     func(int, int, int) (int, int, error)
	FireCommandCallData []struct{ TargetX, TargetY, Enemies int }
}

func (m *IonCannonClientMock) CheckStatus() (*domain.IonCannon, error) {
	callData := struct{}{}
	m.CheckStatusCallData = append(m.CheckStatusCallData, callData)

	if m.CheckStatusFunc != nil {
		return m.CheckStatusFunc()
	}

	return nil, nil
}

func (m *IonCannonClientMock) FireCommand(targetX int, targetY int, enemies int) (casualties int, generation int, err error) {
	callData := struct{ TargetX, TargetY, Enemies int }{targetX, targetY, enemies}
	m.FireCommandCallData = append(m.FireCommandCallData, callData)

	if m.FireCommandFunc != nil {
		return m.FireCommandFunc(targetX, targetY, enemies)
	}

	return 0, 0, nil
}
