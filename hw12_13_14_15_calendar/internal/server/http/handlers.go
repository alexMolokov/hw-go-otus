package internalhttp

import "net/http"

func (s *Server) HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("hello page"))
	}
}
