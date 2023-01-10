package main

import (
	"log"
	"os"
	"os/signal"

	"crg.eti.br/go/mooca/config"
)

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

	log.Println(cfg.Listen)
}
