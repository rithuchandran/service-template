package hotel

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RegionServiceTestSuite struct {
	suite.Suite
	repository *MockRegionRepository
	client     *mockClient
}

func (s *RegionServiceTestSuite) SetupTest() {
	s.repository = &MockRegionRepository{}
	s.client = &mockClient{}
}

func TestRegionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RegionServiceTestSuite))
}

func (s *RegionServiceTestSuite) TestUpdate() {
	service := NewRegionService(s.repository, s.client)
	mockRegions := Regions{"1": Region{Name: "test region", Id: "1", Type: "city"}}

	s.client.On("getRegions").Return(mockRegions,nil)
	s.repository.On("update", mockRegions).Times(1).Return(nil)

	err := service.Update()

	assert.NoError(s.T(), err)
	s.client.AssertExpectations(s.T())
	s.repository.AssertExpectations(s.T())
}

func (s *RegionServiceTestSuite) TestUpdateShouldReturnClientError() {
	service := NewRegionService(s.repository, s.client)

	s.client.On("getRegions").Return(Regions{}, errors.New("client error"))

	err := service.Update()

	assert.EqualError(s.T(), err, "client error")
	s.client.AssertExpectations(s.T())
}

func (s *RegionServiceTestSuite) TestUpdateShouldReturnRepositoryError() {
	service := NewRegionService(s.repository, s.client)

	mockRegions := Regions{}
	s.client.On("getRegions").Return(mockRegions, nil)
	s.repository.On("update", mockRegions).Times(1).Return(errors.New("repository error"))

	err := service.Update()

	assert.EqualError(s.T(), err, "repository error")
	s.client.AssertExpectations(s.T())
}

func (s *RegionServiceTestSuite) TestSearch() {
	service := NewRegionService(s.repository, s.client)
	expectedRegion := Region{Name: "test region", Id: "1", Type: "city"}

	tt := []struct {
		testDescription string
		destination     string
		mockRegion      Region
		mockError       error
		expectedRegion  Region
		assertion       assert.ErrorAssertionFunc
	}{
		{"ShouldReturnRegion", "test region", expectedRegion, nil, expectedRegion, assert.NoError},
		{"ShouldReturnError", "test ", Region{}, errors.New("error"), Region{}, assert.Error},
	}

	for _, tc := range tt {
		s.T().Run(tc.testDescription, func(t *testing.T) {
			s.repository.On("get", tc.destination).Return(tc.mockRegion, tc.mockError)

			returnedRegion, err := service.Search(tc.destination)

			s.repository.AssertExpectations(t)
			assert.Equal(t, tc.expectedRegion, returnedRegion)
			tc.assertion(t, err)
		})

	}
}
