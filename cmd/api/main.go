package main

import (
	"log"

	"github.com/biboyqg/social/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8081"),
	}
	app := &application{
		config: cfg,
	}
	log.Fatal(app.run(app.mount()))
}
