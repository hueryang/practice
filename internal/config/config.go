package config

import (
	"os"
)

// Config holds process-wide settings loaded from the environment.
type Config struct {
	// Addr is the TCP listen address (e.g. ":8080").
	Addr string
	// DBPath is the SQLite database file path.
	DBPath string
}

// Load reads configuration from the environment with defaults suitable for local dev.
//
//   - ADDR: listen address, default ":8080"
//   - DB_PATH: SQLite file path, default "./data/mail.db"
func Load() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/mail.db"
	}
	return Config{Addr: addr, DBPath: dbPath}
}
