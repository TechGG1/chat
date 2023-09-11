package main

import (
	"github.com/TechGG1/chat/server/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}
}
