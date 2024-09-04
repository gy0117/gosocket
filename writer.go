package gosocket

import (
	"bytes"
	"errors"
	"github.com/gy/gosocket/internal/tools"
	"github.com/gy/gosocket/internal/xerr"
	"github.com/gy/gosocket/pkg/bufferpool"
)

func (wsConn *WsConn) WriteMessage(opcode Opcode, payload []byte) error {
	err := wsConn.writeMessage(opcode, payload)
	if err != nil {
		wsConn.handleErrorEvent(err)
	}
	return err
}

func (wsConn *WsConn) WriteString(str string) error {
	return wsConn.WriteMessage(OpcodeTextFrame, tools.StringToBytesStandard(str))
}

func (wsConn *WsConn) WritePing(payload []byte) error {
	return wsConn.WriteMessage(OpcodePingFrame, payload)
}

func (wsConn *WsConn) WritePong(payload []byte) error {
	return wsConn.WriteMessage(OpcodePongFrame, payload)
}

func (wsConn *WsConn) writeMessage(opcode Opcode, payload []byte) error {
	// TODO 状态检查
	wsConn.lock.Lock()
	defer wsConn.lock.Unlock()

	n := len(payload)
	if n > wsConn.config.MaxWritePayloadSize {
		return xerr.NewError(xerr.ErrCloseTooLarge, errors.New("payload size more than MaxWritePayloadSize"))
	}
	if wsConn.config.OpenUTF8Check && !isValidText(opcode, payload) {
		return xerr.NewError(xerr.ErrCloseUnSupported, errors.New("invalid text encode, must be utf-8 encode"))
	}
	frame, err := wsConn.createFrame(opcode, payload)
	if err != nil {
		return err
	}
	_, err = wsConn.conn.Write(frame.Bytes())
	if err != nil {
		return err
	}
	bufferpool.Pools.Put(frame)
	return nil
}

func (wsConn *WsConn) createFrame(opcode Opcode, payload []byte) (*bytes.Buffer, error) {
	n := len(payload)
	buf := bufferpool.Pools.Get(headerFrameLen + n)
	f := &Frame{}
	headerLen, maskingKey := f.CreateHeader(true, opcode, wsConn.server, n)

	buf.Write(f.Header[:headerLen])
	if !wsConn.server {
		unMaskPayload(payload, maskingKey)
	}
	buf.Write(payload)
	return buf, nil
}
