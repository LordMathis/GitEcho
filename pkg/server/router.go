package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRouter(apiHandler *APIHandler) *chi.Mux {
	router := chi.NewRouter()

	staticURLPattern := "/static/*"

	// Set up middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // Add your allowed origins here
	}))

	router.Get(staticURLPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir(apiHandler.StaticDir))).ServeHTTP(w, r)
	}))

	router.Post("/api/v1/repository", apiHandler.HandleCreateBackupRepo)
	router.Get("/api/v1/repository/{name}", apiHandler.HandleGetBackupRepoByName)
	router.Get("/api/v1/repository", apiHandler.HandleGetBackupRepos)
	router.Delete("/api/v1/repository/{name}", apiHandler.HandleDelete)
	router.Get("/", apiHandler.HandleIndex)

	return router
}
