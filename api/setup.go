package api

import (
	"net/http"

	"github.com/charmbracelet/log"
)

type Server struct {
	listenAddr string
}

func NewServer(port string) *Server {
	return &Server{
		listenAddr: ":" + port,
	}
}

func (cfg *apiConfig) initFileServer() http.Handler {
	fileHandler := http.FileServer(http.Dir("."))
	stripPref := http.StripPrefix("/app/", fileHandler)

	return cfg.middlewareMetricInc(stripPref)
}

func (s *Server) Start() error {
	config := NewApiConfig()
	log.Infof("Listening to requests on port %s", s.listenAddr)

	// File handler
	fileHandler := config.initFileServer()
	http.Handle("/app/", fileHandler)

	// Route handlers
	http.HandleFunc("GET /api/healthz", healthCheck)
	http.HandleFunc("POST /api/validate_chirp", handleChirps)

	// Admin routes
	http.HandleFunc("GET /admin/metrics", config.handleReqCount)
	http.HandleFunc("POST /admin/reset", config.handleReqCountReset)

	return http.ListenAndServe(s.listenAddr, nil)
}
