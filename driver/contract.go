package driver

// IOPinner IOPinner
type IOPinner interface {
	Value() (int, error)
	SetValue(int) error
	Close() error
}
