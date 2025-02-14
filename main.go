package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/austinlukaschewski/go-work-in-sports/internal/model"
	"github.com/austinlukaschewski/go-work-in-sports/internal/xlog"
	"github.com/austinlukaschewski/go-work-in-sports/src/api/v1/fangraphs"
	"github.com/austinlukaschewski/go-work-in-sports/src/api/v1/teamworkonline"
	"github.com/austinlukaschewski/go-work-in-sports/src/ui"
	"github.com/austinlukaschewski/go-work-in-sports/src/ui/component"
	"github.com/austinlukaschewski/go-work-in-sports/src/ui/page"
)

var middlewareLogger = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xlog.Info(fmt.Sprintf("[%s] %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

var jobKeywords []string = []string{"software", "engineer", "developer", "application", "web"}

func main() {
	mux := http.NewServeMux()

	// File server
	fs := http.FileServer(http.Dir("./static/public"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ui.Page(page.Index([]string{"Teamworkonline", "Fangraphs"}), w, r)
	})

	mux.HandleFunc("GET /api/widget/{section}", func(w http.ResponseWriter, r *http.Request) {
		section := r.PathValue("section")
		var jobs []model.Job
		var err error = nil

		switch section {
		case "fangraphs":
			jobs, err = fangraphs.GetJobs(jobKeywords, 10)
			break
		case "teamworkonline":
			jobs, err = teamworkonline.GetJobs(jobKeywords, 10)
			break
		default:
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ui.Component(component.JobList(jobs), w, r)
	})

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), middlewareLogger(mux))

		if errors.Is(err, http.ErrServerClosed) {
			xlog.Warn("Server shutting down")
		} else {
			xlog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	<-killSig
}
