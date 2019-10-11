package maportd

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/aquilahkj/maport/internal/conn"
	"github.com/aquilahkj/maport/internal/log"
)

type Tunnel struct {
	// time when the tunnel was opened
	start time.Time

	publishConn conn.Conn

	proxyConn conn.Conn

	log.Logger

	receiveBytes int64

	sendBytes int64

	id int32
}

func NewTunnel(pub conn.Conn, pxy conn.Conn) *Tunnel {
	tunnel := &Tunnel{
		start:       time.Now(),
		publishConn: pub,
		proxyConn:   pxy,
		Logger:      log.NewLogger(),
		id:          rand.Int31(),
	}
	// var a int64
	// var b int64
	// tunnel.receiveBytes = &a
	// tunnel.sendBytes = &b
	tunnel.AddLogPrefix(tunnel.String())
	return tunnel
}

func (t *Tunnel) String() string {
	return fmt.Sprintf("tun:%x", t.id)
}

func (t *Tunnel) Close() {
	t.publishConn.Close()
	t.proxyConn.Close()
	t.Debug("The tunnel closed")
}

func (t *Tunnel) JoinPipe() {
	var wait sync.WaitGroup
	pipe := func(dst conn.Conn, src conn.Conn, written *int64) {
		defer wait.Done()
		size := 4 * 1024
		buf := make([]byte, size)
		var err error
		for {
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw > 0 {
					total := *written + int64(nw)
					*written = total
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				break
			}
		}
		if err != nil {
			src.Warn("Connection closed with error %v", err)
		} else {
			src.Debug("Connection closed")
		}
	}

	wait.Add(2)
	go pipe(t.publishConn, t.proxyConn, &t.receiveBytes)
	go pipe(t.proxyConn, t.publishConn, &t.sendBytes)
	t.Debug("The publish connection %v Joined with proxy connection %s", t.publishConn, t.proxyConn)
	wait.Wait()
}
