package main

import (
	"context"
	"fmt"
	"go-rest-api/api"
	"go-rest-api/cfg"
	"go-rest-api/db"
	"go-rest-api/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-sql-driver/mysql"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func run(ctx context.Context, getEnv func(string, string) string) error {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	env := cfg.InitConfig(getEnv)

	cfg := mysql.Config{
		User:                 env.DBUser,
		Passwd:               env.DBPassword,
		Addr:                 env.DBAddress,
		DBName:               env.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sql := db.NewMySQLStorage(cfg)
	db, err := sql.Init()

	if err != nil {
		log.Fatal(err)
	}

	store := storage.NewDbStorage(db)

	svc := api.NewAPIServer(store)

	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: svc,
	}
	fmt.Println("starting server")

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	doneCh := make(chan struct{})

	shutdownCtx := context.Background()
	shutdownCtx, cancel = context.WithTimeout(shutdownCtx, time.Second*2)
	defer cancel()

	log.Println("starting shutdown...")

	go func() {
		time.Sleep(1 * time.Second)
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing db: %s", err)
		}
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down: %s", err)
		}
		close(doneCh)
	}()

	select {
	case <-shutdownCtx.Done():
		fmt.Fprintf(os.Stderr, "failed to shutdown, timeout")
	case <-doneCh:
		log.Println("shutdown successful")
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, getEnv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
