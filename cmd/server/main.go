package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hueryang/practice/internal/config"
	"github.com/hueryang/practice/internal/httpserver"
	"github.com/hueryang/practice/internal/store"
)

func main() {
	cfg := config.Load()
	db, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Printf("store: %v", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	h := httpserver.New(db)
	log.Printf("server listening addr=%s db=%s", cfg.Addr, cfg.DBPath)
	if err := http.ListenAndServe(cfg.Addr, h); err != nil {
		log.Printf("listen: %v", err)
		os.Exit(1)
	}
}
