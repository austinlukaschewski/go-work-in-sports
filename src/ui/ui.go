package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/austinlukaschewski/go-work-in-sports/src/ui/component"
)

func Component(c templ.Component, w http.ResponseWriter, r *http.Request) {
	if err := c.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Page(c templ.Component, w http.ResponseWriter, r *http.Request) {
	if err := component.Layout(c).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
