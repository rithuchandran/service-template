package route_test

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	. "hotels-service-template/route"
	"io"
	"net/http/httptest"
	"testing"
)

type RouteTestSuite struct {
	suite.Suite
	mockHandler *MockRegionHandler
	router      *Router
	rr          *httptest.ResponseRecorder
}

func (s *RouteTestSuite) SetupSuite() {
	s.mockHandler = &MockRegionHandler{}
}

func (s *RouteTestSuite) SetupTest() {
	s.router = New(mux.NewRouter())
	s.rr = httptest.NewRecorder()
}

func TestRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (s *RouteTestSuite) TestRouting() {
	s.router.Configure(s.mockHandler)

	tt := []struct {
		httpMethod        string
		handlerMethodName string
		targetEndpoint    string
		body              io.Reader
	}{
		{httpMethod: "GET", handlerMethodName: "Update", targetEndpoint: "/update"},
		{httpMethod: "GET", handlerMethodName: "Search", targetEndpoint: "/search"},
	}

	for _, tc := range tt {
		req := httptest.NewRequest(tc.httpMethod, tc.targetEndpoint, tc.body)
		s.mockHandler.On(tc.handlerMethodName, s.rr, mock.AnythingOfType("*http.Request")).Return()
		s.router.ServeHTTP(s.rr, req)
		s.mockHandler.AssertExpectations(s.T())
	}
}

func (s *RouteTestSuite) TestWrap() {
	s.router.Configure(s.mockHandler)
	req := httptest.NewRequest("GET", "/update", nil)
	mw1 := &MockMiddleware{}
	mw2 := &MockMiddleware{}
	s.mockHandler.On("Update", s.rr, mock.AnythingOfType("*http.Request")).Return()
	mw2.On("Do", mock.AnythingOfType("*route.Router")).Once().Return()
	mw1.On("Do", mock.AnythingOfType("http.HandlerFunc")).Once().Return()

	h := s.router.Wrap(mw2.Do, mw1.Do)
	h.ServeHTTP(s.rr, req)

	s.mockHandler.AssertExpectations(s.T())
	mw1.AssertExpectations(s.T())
	mw2.AssertExpectations(s.T())

}
