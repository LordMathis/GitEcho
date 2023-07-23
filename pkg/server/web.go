package server

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (a *APIHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {

	templatePah := filepath.Join(a.templatesDir, "index.html")
	tmpl, err := template.ParseFiles(templatePah)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
