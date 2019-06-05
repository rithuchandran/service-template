package hotel

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

type MockRegionRepository struct {
	mock.Mock
}

func (m *MockRegionRepository) update(regions Regions) error {
	fmt.Println("Mocked get destination names function")
	args := m.Called(regions)
	if args[0] != nil {
		return args[0].(error)
	}
	return nil
}

func (m *MockRegionRepository) get(dest string) (Region, error) {
	fmt.Println("Mocked get destination names function")
	args := m.Called(dest)
	fmt.Println("Args extracted are: ", args[0], args[1])
	if args[1] != nil {
		return args[0].(Region), args[1].(error)
	}
	return args[0].(Region), nil
}
