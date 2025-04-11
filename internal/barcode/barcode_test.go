package barcode

import (
	"testing"

	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateBarcodeString(t *testing.T) {
	b := New(nil) // Logger is not needed for this test

	cID := "C123"
	oID := "O456"
	pID := 789

	expected := "C123*O456*789"
	result := b.CreateBarcodeString(cID, oID, pID)

	assert.Equal(t, expected, result, "Barcode string does not match expected format")
}

func TestResolveBarcode(t *testing.T) {
	b := New(nil) // Logger is not needed for this test

	barcode := "C123*O456*789"
	expected := model.BarcodeFields{
		CustomerID: "C123",
		OrderID:    "O456",
		ProductID:  789,
	}

	result, err := b.ResolveBarcode(barcode)

	assert.NoError(t, err, "Error should not occur while resolving barcode")
	assert.Equal(t, expected, result, "Resolved barcode fields do not match expected values")
}

func TestResolveBarcode_InvalidFormat(t *testing.T) {
	b := New(nil) // Logger is not needed for this test

	barcode := "Invalid*Barcode"
	_, err := b.ResolveBarcode(barcode)

	assert.Error(t, err, "Error should occur for invalid barcode format")
}

func TestResolveBarcode_InvalidProductID(t *testing.T) {
	b := New(nil) // Logger is not needed for this test

	barcode := "C123*O456*InvalidID"
	_, err := b.ResolveBarcode(barcode)

	assert.Error(t, err, "Error should occur for invalid product ID in barcode")
}
