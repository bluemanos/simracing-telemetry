package server

type HandleConnection func(buffer []byte, port int)

type Server interface {
	Run(fn HandleConnection, port int) error
	Close() error
}

// NewServer creates a new Server
func NewServer(addr string) Server {
	return &UDPServer{
		Addr: addr,
	}
}
