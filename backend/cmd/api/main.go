package main

import (
	"log"

	"github.com/AelcioJozias/vibe-invest/backend/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
