package barcode

import "github.com/alp-tahta/warehouse/internal/model"

type Barcoder interface {
	CreateBarcodeString(cID, oID string, pID int) string
	ResolveBarcode(barcode string) (bf model.BarcodeFields, e error)
	// Create(text string) error
	// Read(file *os.File) error
}
