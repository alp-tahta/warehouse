package barcode

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

type Barcode struct {
	l *slog.Logger
}

func NewBarcode(l *slog.Logger) *Barcode {
	return &Barcode{
		l: l,
	}
}

func (b *Barcode) Create(text, fileName string) error {
	// Generate barcode
	code, err := code128.Encode(text)
	if err != nil {
		return err
	}

	scaledCode, err := barcode.Scale(code, 300, 100)
	if err != nil {
		return err
	}

	// Permission: 0777 (rwxrwxrwx)
	perm := os.FileMode(0777)

	// Create file with 0777 permissions
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, scaledCode)
	if err != nil {
		return err
	}
	b.l.Info("✅ Barcode generated:", "info", fileName)
	return nil
}

func (b *Barcode) Read(fileName string) error {
	barcodeFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer barcodeFile.Close()

	img, _, err := image.Decode(barcodeFile)
	if err != nil {
		return err
	}

	source := gozxing.NewLuminanceSourceFromImage(img)
	bitmap, err := gozxing.NewBinaryBitmap(gozxing.NewGlobalHistgramBinarizer(source))
	reader := oned.NewCode128Reader()
	result, err := reader.Decode(bitmap, nil)
	if err != nil {
		return err
	}

	b.l.Info("✅ Decoded barcode content:", "info", result.GetText())
	return nil
}

func CreateBarcodeString(cID, oID string, pID int) string {
	return fmt.Sprintf("%s*%s*%d", cID, oID, pID)
}

func ResolveBarcode(barcode string) (b model.BarcodeFields, e error) {
	parts := strings.Split(barcode, "*")
	if len(parts) != 3 {
		return b, errors.New("Invalid barcode format")
	}

	b.CustomerID = parts[0]
	b.OrderID = parts[1]
	pID, err := strconv.Atoi(parts[2])
	if err != nil {
		return b, err
	}
	b.ProductID = pID

	return b, nil
}
