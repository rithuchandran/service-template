package hotel

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RegionServiceTestSuite struct {
	suite.Suite
	repository *MockRegionRepository
	client     *mockClient
}

func (s *RegionServiceTestSuite) SetupSuite() {
	s.repository = &MockRegionRepository{}
	s.client = &mockClient{}
}

func TestRegionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RegionServiceTestSuite))
}

func (s *RegionServiceTestSuite) TestUpdate() {
	service := NewRegionService(s.repository, s.client)
	expectedRegions := Regions{"1": Region{Name: "test region", Id: "1", Type: "city"}}

	s.client.On("getRegions").Return(expectedRegions)
	s.repository.On("update", expectedRegions).Return(nil)

	service.Update()

	s.client.AssertExpectations(s.T())
	s.repository.AssertExpectations(s.T())
}

func (s *RegionServiceTestSuite) TestSearch() {
	service := NewRegionService(s.repository, s.client)
	expectedRegion := Region{Name: "test region", Id: "1", Type: "city"}

	s.repository.On("get", "test region").Return(expectedRegion,nil)

	service.Search("test region")

	s.repository.AssertExpectations(s.T())
}
