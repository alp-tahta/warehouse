package main

import (
	localbarcode "github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/logger"
	"github.com/alp-tahta/warehouse/internal/model"
)

func main() {
	// Init Logger
	l := logger.Init()

	fileName := "barcode.png"

	barcodeText := localbarcode.CreateBarcodeString(
		model.Order{
			OrderID:    1,
			CustomerID: 1,
			OrderItems: nil,
		},
	)

	// Init barcode
	bc := localbarcode.NewBarcode(l)
	// Create barcode
	err := bc.Create(barcodeText, fileName)
	if err != nil {
		l.Error("Error while creating barcode", "err", err.Error())
	}
	// Read barcode
	err = bc.Read(fileName)
	if err != nil {
		l.Error("Error while reading barcode", "err", err.Error())
	}
}
