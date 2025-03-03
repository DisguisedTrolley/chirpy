package api

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/DisguisedTrolley/chirpy/app/internal/database"
	"github.com/charmbracelet/log"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func NewApiConfig() *apiConfig {
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Errorf("Unable to upen db connection: %s", err)
		return nil
	}

	dbQueries := database.New(db)

	return &apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
	}
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
