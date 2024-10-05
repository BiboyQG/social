package main

import (
	"github.com/biboyqg/social/internal/db"
	"github.com/biboyqg/social/internal/env"
	"github.com/biboyqg/social/internal/store"
	"go.uber.org/zap"
)

//	@title			Social Network API
//	@description	API for a Social Network server.
//	@termsOfService	https://github.com/biboyqg/social/blob/main/TERMS_OF_SERVICE.md

//	@contact.name	API Support
//	@contact.url	https://github.com/biboyqg/social
//	@contact.email	banghao2@illinois.edu

//	@license.name	MIT
//	@license.url	https://github.com/biboyqg/social/blob/main/LICENSE

//	@host		localhost:8081
//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				JWT authorization header
func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8081"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8081"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "dev"),
		version: env.GetString("VERSION", "0.0.1"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
