package hotel_handler

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"hotels-service-template/hotel"
)

type MockRegionService struct {
	mock.Mock
}

func (m *MockRegionService) Update() error {
	fmt.Println("MockRegionService Update method called")
	args := m.Called()
	fmt.Println("args extracted are : ", args)
	if args[0] != nil {
		return args[0].(error)
	}
	return nil
}

func (m *MockRegionService) Search(destination string) (hotel.Region, error) {
	fmt.Println("MockRegionService Search method called")
	args := m.Called(destination)
	fmt.Println("args extracted are : ", args[0])
	if args[1] != nil {
		return args[0].(hotel.Region), args[1].(error)
	}
	return args[0].(hotel.Region), nil
}
