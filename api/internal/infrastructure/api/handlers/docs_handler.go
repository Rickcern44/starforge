package handlers

import (
	"log/slog"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/go-chi/chi/v5"
)

type DocsHandler struct{}

func RegisterDocsEndpoints(r chi.Router, specPath string) {
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		html, err := scalargo.NewV2(
			scalargo.WithSpecDir(specPath),
			scalargo.WithBaseFileName("swagger.yaml"),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	})
	slog.Info("Docs handler registered and available at /docs")
}
