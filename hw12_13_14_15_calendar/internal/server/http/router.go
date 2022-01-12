package internalhttp

import (
	"net/http"
)

func (s *Server) NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", s.HomeHandler())
	return mux
}
