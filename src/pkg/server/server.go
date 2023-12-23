package server

type HandleConnection func(buffer []byte)

type Server interface {
	Run(fn HandleConnection) error
	Close() error
}

// NewServer creates a new Server using given protocol
// and addr.
func NewServer(addr string) Server {
	return &UDPServer{
		Addr: addr,
	}
}
