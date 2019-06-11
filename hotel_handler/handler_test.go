package hotel_handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	. "hotels-service-template/hotel"
	"hotels-service-template/hotel_handler"
	"net/http/httptest"
	"testing"
)

type RegionHandlerTestSuite struct {
	suite.Suite
	service *hotel_handler.MockRegionService
}

func (s *RegionHandlerTestSuite) SetupSuite() {
	s.service = &hotel_handler.MockRegionService{}
}

func TestRegionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(RegionHandlerTestSuite))
}

func (s *RegionHandlerTestSuite) TestUpdate() {
	req := httptest.NewRequest("GET", "/update", nil)
	handler := hotel_handler.NewRegionHandler(s.service)
	expectedErrorResponse := bytes.NewBuffer(nil)
	_ = json.NewEncoder(expectedErrorResponse).Encode(hotel_handler.Error{HttpStatus: 500,
		Message: "db error"})

	tt := []struct {
		testDescription  string
		mockError        error
		expectedResponse *bytes.Buffer
	}{
		{"ShouldNotReturnError", nil, bytes.NewBuffer([]byte("update successful"))},
		{"ShouldReturnError", errors.New("db error"), expectedErrorResponse},
	}

	for _, tc := range tt {
		s.T().Run(tc.testDescription, func(t *testing.T) {
			rr := httptest.NewRecorder()
			s.service.On("Update").Times(1).Return(tc.mockError)
			handler.Update(rr, req)
			s.service.AssertExpectations(t)
			assert.Equal(t, tc.expectedResponse, rr.Body)
		})
	}
}

func (s *RegionHandlerTestSuite) TestSearch() {
	req := httptest.NewRequest("GET", "/search?destination=first", nil)
	handler := hotel_handler.NewRegionHandler(s.service)

	testRegion := Region{Id: "1", Name: "first"}
	expectedRegionResponse := bytes.NewBuffer(nil)
	_ = json.NewEncoder(expectedRegionResponse).Encode(testRegion)
	expectedErrorResponse := bytes.NewBuffer(nil)
	_ = json.NewEncoder(expectedErrorResponse).Encode(hotel_handler.Error{HttpStatus: 500,
		Message: "error"})

	tt := []struct {
		testDescription  string
		mockError        error
		mockRegion       Region
		expectedResponse *bytes.Buffer
	}{
		{"ShouldNotReturnError", nil, testRegion, expectedRegionResponse},
		{"ShouldReturnError", errors.New("error"), Region{}, expectedErrorResponse},
	}
	for _, tc := range tt {
		s.T().Run(tc.testDescription, func(t *testing.T) {
			rr := httptest.NewRecorder()
			s.service.On("Search", "first").Times(1).Return(tc.mockRegion, tc.mockError)
			handler.Search(rr, req)
			s.service.AssertExpectations(t)
			assert.Equal(t, tc.expectedResponse, rr.Body)
		})
	}
}
