package dash

import (
	"net"
)

// A server is an object representing a dash server.
type server struct {
}

// Creates and initializes a server
func NewServer(iface string, port uint) (this *server, err error) {
	//Create server
	this = &server{}

	return this, nil
}

func (this *server) Run() (chan string, chan error) {
	out := make(chan string)
	err := make(chan error)
	go this.run(out, err)
	return out, err
}

func (this *server) run(outchan chan string, errchan chan error) {
	defer close(outchan)
	defer close(errchan)
	ln, err := net.Listen("tcp", "127.0.0.1:8443")
	if err != nil {
		errchan <- err
	} else {

		for {
			conn, err := ln.Accept()
			if err != nil {
				errchan <- err
			}
			outchan <- conn.RemoteAddr().String()
			conn.Close()
		}
	}
	return
}
