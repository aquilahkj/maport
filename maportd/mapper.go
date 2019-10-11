package maportd

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aquilahkj/maport/internal/conn"
	"github.com/aquilahkj/maport/internal/log"
)

// Mapper info
type Mapper struct {
	// time when the mapper was opened
	start time.Time

	// public port
	port int

	// dest address
	dest []string

	// tcp listener
	listener *net.TCPListener

	// logger
	log.Logger

	// closing
	closing int32

	// mapper name
	name string

	tunnels map[string]*Tunnel

	sync.Mutex
}

// Shutdown stop the mapper service
func (t *Mapper) Shutdown() {
	t.Info("Shutting down")

	// mark that we're shutting down
	atomic.StoreInt32(&t.closing, 1)

	// if we have a public listener (this a raw TCP tunnel), shut it down
	if t.listener != nil {
		t.listener.Close()
	}

}

// NewMapper create a new tunnel from a registration message
func NewMapper(port int, destAddr string) (t *Mapper, err error) {
	if port <= 0 || port > 65534 {
		err := fmt.Errorf("bind port %v out of range", port)
		return nil, err
	}
	destArr := strings.Split(destAddr, ",")
	if len(destArr) == 0 {
		err := errors.New("the dest address length is 0")
		return nil, err
	}
	t = &Mapper{
		start:   time.Now(),
		Logger:  log.NewLogger(fmt.Sprintf("map:%v", port)),
		port:    port,
		dest:    destArr,
		tunnels: make(map[string]*Tunnel),
	}

	t.Info("mapper created")

	return t, nil
}

// Start startup the mapper
func (t *Mapper) Start() error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: t.port})
	if err != nil {
		t.Error("Failed to binding TCP listener: %v", err)
		return err
	}
	t.listener = listener
	go t.listenTCP()
	return nil
}

func (t *Mapper) listenTCP() {
	for {
		defer func() {
			if r := recover(); r != nil {
				t.Warn("Failed to listen tcp port %v: %v", t.port, r)
			}
		}()

		// accept public connections
		tcpConn, err := t.listener.AcceptTCP()

		if err != nil {
			// not an error, we're shutting down this tunnel
			if atomic.LoadInt32(&t.closing) == 1 {
				return
			}

			t.Error("Failed to accept new TCP connection: %v", err)
			continue
		}

		go t.handleAccept(tcpConn)
	}
}

func (t *Mapper) handleAccept(tcpConn net.Conn) {
	pubConn := conn.NewConn(tcpConn, "pub")
	pubConn.Debug("Open new connection %v", pubConn.RemoteAddr())
	var pxyConn conn.Conn
	defer func() {
		if r := recover(); r != nil {
			pubConn.Warn("Occur error in handleAccept: %v", r)
		}
	}()

	var dest string
	if len(t.dest) > 1 {
		i := rand.Int() % len(t.dest)
		dest = t.dest[i]
	} else {
		dest = t.dest[0]
	}

	dialConn, err := net.Dial("tcp", dest)
	if err != nil {
		t.Error("Dial dest address %v failed: %v", dest, err)
		pubConn.Close()
		return
	}

	pxyConn = conn.NewConn(dialConn, "pxy")
	pxyConn.Debug("Dial new connection %v", pxyConn.RemoteAddr())
	pxyConn.SetDeadline(time.Time{})

	tunnel := NewTunnel(pubConn, pxyConn)
	tunnel.Debug("Create a tunnel between %v and %v", pubConn.FullInfo(), pxyConn.FullInfo())
	t.insertTunnel(tunnel)
	tunnel.JoinPipe()
	t.releaseTunnel(tunnel)
	tunnel.Close()
}

func (t *Mapper) insertTunnel(tunnel *Tunnel) {
	defer t.Unlock()
	t.Lock()
	t.tunnels[tunnel.String()] = tunnel
}

func (t *Mapper) releaseTunnel(tunnel *Tunnel) {
	defer t.Unlock()
	t.Lock()
	delete(t.tunnels, tunnel.String())
}
