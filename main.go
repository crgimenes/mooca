package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"crg.eti.br/go/mooca/config"
	"crg.eti.br/go/mooca/session"
)

var (
	sc *session.Control

	//go:embed assets
	assets embed.FS
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := sc.Get(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// renew session
	sc.Save(w, sid, sd)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt)
		<-sc

		log.Println("shutting down...")

		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: check if database is up
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		w.Write([]byte(`{"status": "ok"}`))

	})
	mux.Handle("/assets/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(fmt.Sprintf("url: %s", r.URL.Path)))
	})

	mux.HandleFunc("/login/", loginHandler)
	mux.HandleFunc("/", homeHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on port %d\n", cfg.Port)
	log.Fatal(s.ListenAndServe())
}
