package route

import (
	"fmt"
	"net/http"
)

type MockHandler struct {
	Request *http.Request
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Request = r
	fmt.Println("Mock Handler Func called")
}