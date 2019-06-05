package route

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockRegionHandler struct {
	mock.Mock
}

func (m *MockRegionHandler) Search(w http.ResponseWriter, r *http.Request){
	fmt.Println("mockRegionHandler search method called")
	m.Called(w, r)
}

func (m *MockRegionHandler) Update(w http.ResponseWriter, r *http.Request){
	fmt.Println("mockRegionHandler update method called")
	m.Called(w, r)
}

