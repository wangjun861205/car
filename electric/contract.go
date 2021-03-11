package electric

// IOPinner IOPinner
type IOPinner interface {
	SetValue(value int) error
	Value() (int, error)
	Close() error
}
