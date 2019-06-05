package route_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"hotels-service-template/route"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ContentTypeSetterTestSuite struct {
	suite.Suite
	mockHandler *route.MockHandler
}

func (s *ContentTypeSetterTestSuite) SetupTest() {
	s.mockHandler = &route.MockHandler{}
}

func TestContentTypeSetterTestSuite(t *testing.T) {
	suite.Run(t, new(ContentTypeSetterTestSuite))
}

func (s *ContentTypeSetterTestSuite) TestSetContentTypeHeader() {
	expectedContentType := []string{"application/json; charset=utf-8"}
	req, _ := http.NewRequest("GET", "/setContentType", nil)
	rr := httptest.NewRecorder()

	route.SetContentTypeHeader(s.mockHandler).ServeHTTP(rr, req)

	obtainedResponseHeaders := rr.Header()
	obtainedContentType, ok := obtainedResponseHeaders["Content-Type"]
	assert.True(s.T(), ok)
	assert.Equal(s.T(), expectedContentType, obtainedContentType)

}
