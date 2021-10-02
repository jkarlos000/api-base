package main

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/errors"
	"backend/internal/healthcheck"
	"backend/internal/user"
	"backend/pkg/accesslog"
	"backend/pkg/dbcontext"
	"backend/pkg/log"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"

	f "github.com/go-ozzo/ozzo-routing/v2/file"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// check if path ssl exists
	/*if path, err := os.Getwd(); err == nil {
		if err := os.Mkdir(path+"/certs", 0755); !os.IsExist(err) {
			logger.Info("Creating path: " + path + "/certs")
		}
	}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("docker-core.ml"), //Your domain here
		Cache:      autocert.DirCache("certs"),               //Folder for storing certificates
	}

	if path, err := os.Getwd(); err == nil {
		if err := os.Mkdir(path+"/storage", 0755); !os.IsExist(err) {
			if err := os.Mkdir(path+"/storage/diagrams", 0755); !os.IsExist(err) {
				logger.Info("Creating path: " + path + "/storage/diagrams")
			}
		}
	}*/

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to the database
	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), cfg),
		/*TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},*/
	}

	// go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	/*if err := hs.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}*/
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
	if _, err = os.Stat("temp"); os.IsNotExist(err) {
		_ = os.Mkdir("temp", 0777)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	authHandler := auth.Handler(cfg.JWTSigningKey)

	// lógica para backend.

	/*album.RegisterHandlers(rg.Group(""),
		album.NewService(album.NewRepository(db, logger), logger),
		authHandler, logger,
	)*/

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(db, cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	user.RegisterHandlers(rg.Group(""),
		user.NewService(user.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	router.Get("/*", f.Server(f.PathMap{
		"/v1/diagrams": "/storage/diagrams",
	}))

	return router
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
