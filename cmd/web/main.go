package main

import (
	"context"
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

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.guimochila.com/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
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

	// Initialize a new intance of our application struct
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", slog.String("addr", ":4000"))

	svr := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	svrErrors := make(chan error, 1)
	go func() {
		svrErrors <- svr.ListenAndServe()
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
