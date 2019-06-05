package route

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) Do(next http.Handler) http.Handler {
	m.Called(next)
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(responseWriter, request)
	})
}
