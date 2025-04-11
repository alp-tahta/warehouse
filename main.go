package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/config"
	"github.com/alp-tahta/warehouse/internal/handler"
	"github.com/alp-tahta/warehouse/internal/logger"
	"github.com/alp-tahta/warehouse/internal/repository"
	"github.com/alp-tahta/warehouse/internal/routes"
	"github.com/alp-tahta/warehouse/internal/service"
)

func main() {
	// Init Logger
	l := logger.Init()

	config, err := config.BuiltConfig()
	if err != nil {
		l.Error("Error while creating config", "error", err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Error("Error while connecting DB", "error", err)
		os.Exit(1)
	}

	// Ensure the connection is available
	if err = db.Ping(); err != nil {
		l.Error("Error while pinging DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	mux := http.NewServeMux()

	barcode := barcode.New(l)
	repository := repository.New(l, barcode, db)
	service := service.New(l, barcode, repository)
	handler := handler.New(l, service)

	routes.RegisterRoutes(mux, handler)

	srv := http.Server{
		Addr:              config.Port,
		Handler:           mux,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	// Create a channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		l.Info("Server starting at", "port", config.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		l.Error("Failed to start server", "error", err)
		os.Exit(1)

	case sig := <-shutdown:
		l.Info("Start shutdown", "signal", sig)

		// Give outstanding requests 5 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			l.Error("Graceful shutdown did not complete", "error", err)
			if err := srv.Close(); err != nil {
				l.Error("Could not stop server", "error", err)
			}
		}
	}

	// barcodeText := localbarcode.CreateBarcodeString("1", "1", 1)

	// // Init barcode
	// bc := localbarcode.NewBarcode(l)
	// // Create barcode
	// err = bc.Create(barcodeText, barcodeText)
	// if err != nil {
	// 	l.Error("Error while creating barcode", "err", err.Error())
	// 	//os.Exit(1)
	// }
	// // Read barcode
	// err = bc.Read(barcodeText)
	// if err != nil {
	// 	l.Error("Error while reading barcode", "err", err.Error())
	// 	//os.Exit(1)
	// }
}
