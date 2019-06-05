package hotel_handler

import (
	"encoding/json"
	"fmt"
	"hotels-service-template/hotel"
	"net/http"
)

type RegionHandlerInt interface {
	Search(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type RegionHandler struct {
	service hotel.RegionServiceInt
}

type Error struct {
	HttpStatus int
	Message    string
}

func NewRegionHandler(regionService hotel.RegionServiceInt) *RegionHandler {
	return &RegionHandler{
		service: regionService,
	}
}

func (h *RegionHandler) Search(w http.ResponseWriter, r *http.Request) {
	destination := r.URL.Query().Get("destination")
	region, err := h.service.Search(destination)
	if err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(region)
}

func (h *RegionHandler) Update(w http.ResponseWriter, r *http.Request) {
	err := h.service.Update()
	if err != nil {
		fmt.Println("***********************************************")
		handleError(err, w, http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "update successful")
}

func handleError(err error, writer http.ResponseWriter, httpStatusCode int) {
	writer.WriteHeader(httpStatusCode)
	_ = json.NewEncoder(writer).Encode(Error{Message: err.Error(), HttpStatus: httpStatusCode})

}
