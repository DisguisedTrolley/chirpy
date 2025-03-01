package api

import (
	"fmt"
	"net/http"
)

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handleReqCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) handleReqCountReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
