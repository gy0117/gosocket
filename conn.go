package gows

import (
	"bufio"
	"net"
)

// WsConn websocket connection
type WsConn struct {
	// Underlying network connection
	conn net.Conn

	bufReader *bufio.Reader
}
