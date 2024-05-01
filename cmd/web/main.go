package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.guimochila.com/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "", "MySQL data source name")
	flag.Parse()

	// Start logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	if *dsn == "" {
		logger.Error("-dsn flag is missing or can't be empty")
		os.Exit(1)
	}

	db, err := openDB(*dsn)
	fmt.Println(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	//Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance
	formdDecoder := form.NewDecoder()

	// Initialize a new Session Manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize a new intance of our application struct
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formdDecoder,
		sessionManager: sessionManager,
	}

	logger.Info("starting server", slog.String("addr", ":4000"))

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	svr := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	svrErrors := make(chan error, 1)
	go func() {
		svrErrors <- svr.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	shutdown := gracefulShutdown(svr)

	select {
	case err := <-svrErrors:
		shutdown(err)
	case sig := <-quit:
		shutdown(sig)
	}
	logger.Error("Server exiting")
}

func gracefulShutdown(svr *http.Server) func(reason interface{}) {
	return func(reason interface{}) {
		log.Println("Server Shutdown:", reason)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			log.Println("Error Gracefully Shutting Down API:", err)
		}
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
