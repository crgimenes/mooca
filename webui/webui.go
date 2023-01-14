package webui

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

func handleIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))

	f, err := assets.Open("assets/favicon.ico")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: check if database is up
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.Write([]byte(`{"status": "ok"}`))

}

func Mux(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("/assets/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("/favicon.ico", handleIcon)
	mux.HandleFunc("/healthcheck/", healthcheckHandler)
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(fmt.Sprintf("url: %s", r.URL.Path)))
	})

	mux.HandleFunc("/login/", loginHandler)
	mux.HandleFunc("/", homeHandler)

	return mux
}
