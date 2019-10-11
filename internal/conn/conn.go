package conn

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/aquilahkj/maport/internal/log"
)

// Conn the connection interface
type Conn interface {
	net.Conn
	log.Logger
	FullInfo() string
}

type mapConn struct {
	net.Conn
	log.Logger
	id     int32
	prefix string
}

func (c *mapConn) String() string {
	return fmt.Sprintf("%s:%x", c.prefix, c.id)
}

func (c *mapConn) FullInfo() string {
	return fmt.Sprintf("%s:%x(%v<->%v)", c.prefix, c.id, c.LocalAddr(), c.RemoteAddr())
}

func NewConn(conn net.Conn, prefix string) Conn {
	wrapConn := &mapConn{conn, log.NewLogger(), rand.Int31(), prefix}
	wrapConn.AddLogPrefix(wrapConn.String())
	return wrapConn
}
