package route

import (
	"github.com/gorilla/mux"
	"hotels-service-template/hotel_handler"
	"net/http"
)

type Router struct {
	*mux.Router
}

func New(router *mux.Router) *Router {
	return &Router{router}
}

func (r Router) Configure(handler hotel_handler.RegionHandlerInt) {
	r.Handle("/", http.FileServer(http.Dir(".")))
	r.HandleFunc("/search", handler.Search)
	r.HandleFunc("/update", handler.Update)
}

func (r *Router) Wrap(middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var wrappedHandler http.Handler = r
	for _, mw := range middlewares {
		wrappedHandler = mw(wrappedHandler)
	}
	return wrappedHandler
}
