package main

import (
	"os"

	"github.com/sarunask/triviadb-gui/internal/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Setting log filters to filter messages
	// Initialize logger
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	// WEB
	err := web.Run()
	if err != nil {
		log.Fatalf("web error: %v", err)
	}
}
