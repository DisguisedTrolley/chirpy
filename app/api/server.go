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

func (s *Server) Start(cfg *apiConfig) error {
	log.Infof("Listening to requests on port %s", s.listenAddr)

	// File handler
	fileHandler := cfg.initFileServer()
	http.Handle("/app/", fileHandler)

	// Route handlers
	http.HandleFunc("GET /api/healthz", healthCheck)
	http.HandleFunc("POST /api/users", cfg.createUser)
	http.HandleFunc("POST /api/chirps", cfg.addChirp)
	http.HandleFunc("GET /api/chirps", cfg.getChirps)
	http.HandleFunc("GET /api/chirps/{chirpID}", cfg.getSingleChirp)

	// Auth handlers
	http.HandleFunc("POST /api/login", cfg.loginUser)

	// Admin routes
	http.HandleFunc("GET /admin/metrics", cfg.getReqCount)
	http.HandleFunc("POST /admin/reset", cfg.reset)

	return http.ListenAndServe(s.listenAddr, nil)
}
