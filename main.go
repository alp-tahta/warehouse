package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	localbarcode "github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/config"
	"github.com/alp-tahta/warehouse/internal/logger"
	"github.com/alp-tahta/warehouse/internal/model"
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

	barcodeText := localbarcode.CreateBarcodeString(
		model.Order{
			ID:         1,
			CustomerID: 1,
			OrderItems: nil,
		},
	)

	// Init barcode
	bc := localbarcode.NewBarcode(l)
	// Create barcode
	err = bc.Create(barcodeText, barcodeText)
	if err != nil {
		l.Error("Error while creating barcode", "err", err.Error())
		//os.Exit(1)
	}
	// Read barcode
	err = bc.Read(barcodeText)
	if err != nil {
		l.Error("Error while reading barcode", "err", err.Error())
		//os.Exit(1)
	}

	time.Sleep(100 * time.Second)
}
