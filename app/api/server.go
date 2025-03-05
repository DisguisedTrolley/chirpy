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
	http.HandleFunc("PUT /api/users", cfg.updateUser)
	http.HandleFunc("POST /api/chirps", cfg.addChirp)
	http.HandleFunc("GET /api/chirps", cfg.getChirps)
	http.HandleFunc("GET /api/chirps/{chirpID}", cfg.getSingleChirp)
	http.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	// Auth handlers
	http.HandleFunc("POST /api/login", cfg.loginUser)
	http.HandleFunc("POST /api/refresh", cfg.refreshToken)
	http.HandleFunc("POST /api/revoke", cfg.revokeRefreshToken)

	// Admin routes
	http.HandleFunc("GET /admin/metrics", cfg.getReqCount)
	http.HandleFunc("POST /admin/reset", cfg.reset)

	return http.ListenAndServe(s.listenAddr, nil)
}
