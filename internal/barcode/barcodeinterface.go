package barcode

import "os"

type Barcoder interface {
	Create(text string) error
	Read(file *os.File) error
}
