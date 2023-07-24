package server

import (
	"net/http"
	"path/filepath"

	_ "github.com/LordMathis/GitEcho/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			GitEcho API
//	@version		1.0
//	@description	REST API for GitEcho, a tool for backing up Git repositories

//	@license.name	MIT
//	@license.url	http://www.opensource.org/licenses/MIT

//	@BasePath	/api/v1

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
		http.StripPrefix("/static/", http.FileServer(http.Dir(apiHandler.staticDir))).ServeHTTP(w, r)
	}))

	apiRouter := chi.NewRouter()
	apiRouter.Route("/", func(r chi.Router) {
		r.Route("/repository", func(r chi.Router) {
			r.Post("/", apiHandler.HandleCreateBackupRepo)
			r.Get("/", apiHandler.HandleGetBackupRepos)
			r.Route("/{repo_name}", func(r chi.Router) {
				r.Get("/", apiHandler.HandleGetBackupRepoByName)
				r.Delete("/", apiHandler.HandleDeleteBackupRepo)
				r.Route("/storage/", func(r chi.Router) {
					r.Get("/", apiHandler.HandleGetBackupRepoStorages)
					r.Route("/{storage_name}", func(r chi.Router) {
						r.Post("/", apiHandler.HandleAddBackupRepoStorage)
						r.Delete("/", apiHandler.HandleRemoveBackupRepoStorage)
					})
				})
			})
		})

		r.Route("/storage", func(r chi.Router) {
			r.Post("/", apiHandler.HandleCreateStorage)
			r.Get("/", apiHandler.HandleGetStorages)
			r.Route("/{storage_name}", func(r chi.Router) {
				r.Get("/", apiHandler.HandleGetStorageByName)
				r.Delete("/", apiHandler.HandleDeleteStorage)
			})
		})
	})

	router.Get("/", apiHandler.HandleIndex)
	router.Mount("/api/v1", apiRouter)

	filepath.Join(apiHandler.templatesDir, "..", "docs", "swagger.json")

	router.Get("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(apiHandler.templatesDir, "..", "docs", "swagger.json"))
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.json"),
	))

	return router
}
