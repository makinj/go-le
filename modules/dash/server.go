package dash

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// A server is an object representing a dash server.
type server struct {
	lastTriggered map[string]time.Time
}

// Creates and initializes a server
func NewServer(iface string, port uint) (this *server, err error) {
	lt := make(map[string]time.Time)
	//Create server
	this = &server{lastTriggered: lt}

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
			this.handleConnection(conn, outchan, errchan)
		}
	}
	return
}

func (this *server) handleConnection(conn net.Conn, outchan chan string, errchan chan error) {
	ts := time.Now()
	raddr := conn.RemoteAddr().String()
	rhost := strings.Split(raddr, ":")[0]
	last, found := this.lastTriggered[rhost]
	if found && ts.Sub(last).Nanoseconds() > 0 {
		fmt.Println(ts.Sub(last).Nanoseconds())
		outchan <- rhost
	}
	this.lastTriggered[rhost] = ts
	conn.Close()
	return
}
