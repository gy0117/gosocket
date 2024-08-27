package gows

import (
	"bufio"
	"github.com/gy/gosocket/internal"
	"net"
	"net/http"
)

type Upgrade struct {
	options *ServerOptions
}

func NewUpgrade(options *ServerOptions) *Upgrade {
	return &Upgrade{
		options: options,
	}
}

// Upgrade
// 升级HTTP连接成websocket
func (up *Upgrade) Upgrade(w http.ResponseWriter, r *http.Request) (*WsConn, error) {
	return up.upgradeInner(w, r)
}

// 劫持HTTP连接，升级成websocket
func (up *Upgrade) upgradeInner(w http.ResponseWriter, r *http.Request) (*WsConn, error) {
	// 1. 劫持
	netConn, _, err := up.hijack(w)
	if err != nil {
		return nil, err
	}
	// 维护缓冲区池子，不使用hijack返回的reader
	reader := up.options.bufReaderPool.Get()
	reader.Reset(netConn)

	// 2. 升级成websocket
	// 2.1 检查是否符合websocket协议规范
	if err = checkHeader(r); err != nil {
		return nil, err
	}
	// 2.2 构造
	wsConn := &WsConn{
		conn:      netConn,
		bufReader: reader,
	}
	return wsConn, nil
}

func checkHeader(r *http.Request) error {
	if r.Method != http.MethodGet {
		return internal.ErrHandShake
	}
	if r.Header.Get(internal.ConnectionPair.Key) != internal.ConnectionPair.Value {
		return internal.ErrHandShake
	}
	if r.Header.Get(internal.UpgradePair.Key) != internal.UpgradePair.Value {
		return internal.ErrHandShake
	}
	if r.Header.Get(internal.SecWebSocketVersionPair.Key) != internal.SecWebSocketVersionPair.Value {
		return internal.ErrHandShake
	}
	return nil
}

func (up *Upgrade) hijack(w http.ResponseWriter) (net.Conn, *bufio.Reader, error) {
	hi, ok := w.(http.Hijacker)
	if !ok {
		return nil, nil, internal.ErrInternalServer
	}
	netConn, rw, err := hi.Hijack()
	if err != nil {
		return nil, nil, err
	}
	return netConn, rw.Reader, nil
}
