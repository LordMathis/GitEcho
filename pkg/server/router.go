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

	apiRouter := chi.NewRouter()
	apiRouter.Route("/", func(r chi.Router) {
		r.Route("/repository", func(r chi.Router) {
			r.Post("/", apiHandler.HandleCreateBackupRepo)
			r.Get("/", apiHandler.HandleGetBackupRepos)
			r.Get("/{name}", apiHandler.HandleGetBackupRepoByName)
			r.Delete("/{name}", apiHandler.HandleDeleteBackupRepo)
		})

		r.Route("/storage", func(r chi.Router) {
			r.Post("/", apiHandler.HandleCreateStorage)
			r.Get("/", apiHandler.HandleGetStorages)
			r.Get("/{name}", apiHandler.HandleGetStorageByName)
			r.Delete("/{name}", apiHandler.HandleDeleteStorage)
		})
	})

	router.Get("/", apiHandler.HandleIndex)
	router.Mount("/api/v1", apiRouter)

	return router
}
