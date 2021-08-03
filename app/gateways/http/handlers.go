package server

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	promhttp.Handler().ServeHTTP(w, r)
}

func VersionHandler(commit, time string) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		b, err := json.Marshal(struct {
			GitCommitHash string `json:"git_hash"`
			BuildTime     string `json:"time"`
		}{
			GitCommitHash: commit,
			BuildTime:     time,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal version")
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(b)
		if err != nil {
			log.Error().Err(err).Msg("failed to write version body")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
