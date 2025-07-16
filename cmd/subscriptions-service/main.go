package main

import (
	"log"

	"github.com/Yagshymyradov/subscriptions-service/internal/config"
)

func main(){
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Starting service on port %s", cfg.HTTP.Port)
}