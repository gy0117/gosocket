package gosocket

import (
	"bufio"
	"bytes"
	"testing"
	"unsafe"
)

func BenchmarkReadMessage(b *testing.B) {
	b.Run("no compression", func(b *testing.B) {
		upgrade := NewUpgrade(&MockEventHandler{}, nil)
		client := &WsConn{
			conn:         &mockConn{},
			bufReader:    upgrade.options.readerBufPool.Get(),
			eventHandler: upgrade.eventHandler,
			frame:        NewFrame(),
			config:       upgrade.options.CreateConfig(),
			server:       false,
		}

		frame, _ := client.createFrame(OpcodeTextFrame, mockData)

		reader := bytes.NewBuffer(frame.Bytes())

		server := &WsConn{
			conn:         &mockConn{},
			bufReader:    bufio.NewReader(reader),
			eventHandler: upgrade.eventHandler,
			frame:        NewFrame(),
			config:       upgrade.options.CreateConfig(),
			server:       true,
		}

		for i := 0; i < b.N; i++ {
			BufferReset(reader, frame.Bytes())
			server.bufReader.Reset(reader)
			_ = server.readMessage()
		}

	})
}

// 重置buffer底层切片
func BufferReset(b *bytes.Buffer, p []byte) {
	*(*[]byte)(unsafe.Pointer(b)) = p
}
