package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"crg.eti.br/go/mooca/config"
	"crg.eti.br/go/mooca/session"
	"github.com/gorilla/mux"
)

var (
	sc *session.Control

	//go:embed assets
	assets embed.FS
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sid, sd, ok := sc.Get(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// renew session
	sc.Save(w, sid, sd)

	http.Redirect(w, r, "/payments", http.StatusFound)
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

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)

	r.PathPrefix("/assets/").Handler(http.FileServer(http.FS(assets)))

	log.Printf("Listening on port %d\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))

}
