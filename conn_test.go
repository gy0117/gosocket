package gosocket

import (
	"net"
	"testing"
)

func TestConnRecycle(t *testing.T) {
	server, _ := net.Pipe()
	options := new(ServerOptions)
	initServerOptions(options)
	wsConn := &WsConn{
		conn:      server,
		bufReader: options.readerBufPool.Get(),
	}
	wsConn.Recycle = func() {
		wsConn.bufReader.Reset(nil)
		options.readerBufPool.Put(wsConn.bufReader)
		wsConn.bufReader = nil
	}
	wsConn.Recycle()
	t.Log(wsConn.bufReader)
}
