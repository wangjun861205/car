package camera

type client interface {
	Run()
	Stop()
	RegisterHandler(func(b []byte))
	Write(b []byte) error
}

type server interface {
	Run()
	Stop()
	RegisterHandler(func(b []byte))
}
