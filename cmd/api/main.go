package main

import (
	"time"

	"github.com/biboyqg/social/internal/auth"
	"github.com/biboyqg/social/internal/db"
	"github.com/biboyqg/social/internal/env"
	"github.com/biboyqg/social/internal/mailer"
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

//	@host		localhost:8080
//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8081"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8081"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "dev"),
		version: env.GetString("VERSION", "0.0.1"),
		mail: mailConfig{
			exp: env.GetDuration("MAIL_EXP", 3*24*time.Hour),
			gomail: gomailConfig{
				host:     env.GetString("MAIL_HOST", "smtp.gmail.com"),
				port:     env.GetInt("MAIL_PORT", 587),
				username: env.GetString("MAIL_USERNAME", "banghao.ch@gmail.com"),
				password: env.GetString("MAIL_PASSWORD", "password"),
				sender:   env.GetString("MAIL_SENDER", "banghao.ch@gmail.com"),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("JWT_SECRET", "secret"),
				exp:    env.GetDuration("JWT_EXP", 7*24*time.Hour),
				aud:    env.GetString("JWT_AUD", "social"),
				iss:    env.GetString("JWT_ISS", "social"),
			},
		},
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

	mailer := mailer.NewGomailer(
		cfg.mail.gomail.host,
		cfg.mail.gomail.port,
		cfg.mail.gomail.username,
		cfg.mail.gomail.password,
		cfg.mail.gomail.sender,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
