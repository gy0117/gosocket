package gosocket

import (
	_ "embed"
	"github.com/gy/gosocket/pkg/bufferpool"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestCreateFrame(t *testing.T) {
	t.Run("c1", func(t *testing.T) {
		wsConn := WsConn{
			server: true,
		}
		frame, err := wsConn.createFrame(OpcodeTextFrame, []byte("abc"))
		t.Log("frame", string(frame.Bytes()))
		assert.NoError(t, err)
	})
}

//go:embed assets/mock.json
var mockData []byte

func BenchmarkWriteMessage(b *testing.B) {
	b.Run("no compression", func(b *testing.B) {
		upgrade := NewUpgrade(&MockEventHandler{}, nil)
		conn := &WsConn{
			conn:         &mockConn{},
			bufReader:    upgrade.options.readerBufPool.Get(),
			eventHandler: upgrade.eventHandler,
			frame:        NewFrame(),
			config:       upgrade.options.CreateConfig(),
			server:       true,
		}
		for i := 0; i < b.N; i++ {
			_ = conn.WriteMessage(OpcodeTextFrame, mockData)
		}
	})
}

type mockConn struct {
	net.TCPConn
}

func (m mockConn) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var _ EventHandler = (*MockEventHandler)(nil)

type MockEventHandler struct {
}

func (m *MockEventHandler) OnStart(conn *WsConn) {

}

func (m *MockEventHandler) OnPing(conn *WsConn, payload []byte) {

}

func (m *MockEventHandler) OnPong(conn *WsConn, payload []byte) {

}

func (m *MockEventHandler) OnMessage(conn *WsConn, msg *Message) {
	// 写入消息后，将msg回收

	// msg的回收
	bufferpool.Pools.Put(msg.Content)
	msg.Content = nil
}

func (m *MockEventHandler) OnStop(conn *WsConn, err error) {

}
