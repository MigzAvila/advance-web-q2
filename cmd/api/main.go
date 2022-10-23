//Filename: cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"todoapi.miguelavila.net/internals/data"
)

// App Version
const version = "1.0.0"

// App config
type config struct {
	port int
	env  string // dev, stg, prd, etc...1
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// dependencies injections
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config
	//read in the flag that are needed to populate the config ~ flag for using as extra cmd
	flag.IntVar(&cfg.port, "port", 4000, "API port")
	flag.StringVar(&cfg.env, "env", "dev", "(dev | stg | prd)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("TODO_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-open-conns", 25, "PostgreSQL max idle open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-open-time", "15m", "PostgreSQL max connections idle time")
	flag.Parse()

	//create a logger ~ use := for undeclared var
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//create the connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	// log successful connection
	logger.Printf("database connection pool established")

	//create instances of out api
	app := &application{
		config: cfg,
		logger: logger,
		models: *data.NewModels(db),
	}

	//create our http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server at %s", cfg.env, srv.Addr)
	//start the server
	err = srv.ListenAndServe()
	logger.Fatal(err)

}

// openDB return a *sql.DB instance
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// create a context with a 5 section timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
