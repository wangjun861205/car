package camera

import "gocv.io/x/gocv"

// Image image
type Image struct {
	Rows int
	Cols int
	Type gocv.MatType
	Data []byte
}
