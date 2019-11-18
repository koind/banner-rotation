package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

// HttpServer
type HttpServer struct {
	domain string
	router http.Handler
	s      *RotationService
}

// Start fires up the http server
func (s *HttpServer) Start() error {
	return http.ListenAndServe(s.domain, s.router)
}

// NewHTTPServer returns http server that wraps rotation business logic
func NewHTTPServer(handleService *RotationService, domain string) *HttpServer {

	r := mux.NewRouter()
	hs := HttpServer{router: r, domain: domain, s: handleService}

	r.HandleFunc("/banner/add", handleService.AddBannerHandle).Methods("POST")
	r.HandleFunc("/banner/set-transition", handleService.SetTransitionHandle).Methods("POST")
	r.HandleFunc("/banner/select", handleService.SelectBannerHandle).Methods("POST")
	r.HandleFunc("/banner/remove/{id}", handleService.RemoveBannerHandle).Methods("DELETE")

	http.Handle("/", r)

	return &hs
}
