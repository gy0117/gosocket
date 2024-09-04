package gosocket

import (
	"bytes"
	"github.com/gy/gosocket/internal/xerr"
	"log"
	"sync"
)

type BroadcastManager struct {
	opcode  Opcode
	payload []byte
	once    sync.Once
}

func NewBroadcastManager(opcode Opcode, payload []byte) *BroadcastManager {
	return &BroadcastManager{
		opcode:  opcode,
		payload: payload,
	}
}

// Broadcast 向客户端发送广播消息
func (bm *BroadcastManager) Broadcast(conn *WsConn) error {
	var frame *bytes.Buffer
	var err error
	bm.once.Do(func() {
		frame, err = conn.createFrame(true, bm.opcode, bm.payload)
	})

	if err != nil {
		return err
	}
	conn.taskQueue.Push(func() {
		// 具体执行
		if err := bm.writeFrame(conn, frame); err != nil {
			conn.handleErrorEvent(err)
		}
	})
	conn.taskQueue.Execute()
	return nil
}

func (bm *BroadcastManager) writeFrame(conn *WsConn, frame *bytes.Buffer) error {
	if conn.isClose == 1 {
		return xerr.NewError(xerr.ErrConnClosed, nil)
	}
	conn.lock.Lock()
	defer conn.lock.Unlock()
	_, err := conn.conn.Write(frame.Bytes())
	return err
}

func (bm *BroadcastManager) stop() {
	log.Println("broadcast stop...")
}
