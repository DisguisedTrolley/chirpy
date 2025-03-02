package api

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/DisguisedTrolley/chirpy/internal/database"
	"github.com/charmbracelet/log"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func NewApiConfig() *apiConfig {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Errorf("Unable to upen db connection: %s", err)
		return nil
	}

	dbQueries := database.New(db)

	return &apiConfig{
		dbQueries: dbQueries,
	}
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
