package gosocket

import (
	"bytes"
	"github.com/gy/gosocket/internal"
)

func (wsConn *WsConn) WriteMessage(opcode Opcode, payload []byte) error {
	err := wsConn.writeMessage(opcode, payload)
	wsConn.handleErrorEvent(err)
	return err
}

func (wsConn *WsConn) WriteString(str string) error {
	return wsConn.WriteMessage(OpcodeTextFrame, StringToBytesStandard(str))
}

func (wsConn *WsConn) WritePing(payload []byte) error {
	return wsConn.WriteMessage(OpcodePingFrame, payload)
}

func (wsConn *WsConn) WritePong(payload []byte) error {
	return wsConn.WriteMessage(OpcodePongFrame, payload)
}

func (wsConn *WsConn) writeMessage(opcode Opcode, payload []byte) error {
	// TODO 状态检查

	n := len(payload)
	if n > wsConn.config.MaxWritePayloadSize {
		return internal.ErrCloseTooLarge
	}
	if wsConn.config.OpenUTF8Check && !isValidText(opcode, payload) {
		return internal.NewXError(internal.ErrCloseUnSupported, internal.ErrTextEncode)
	}
	frame, err := wsConn.createFrame(opcode, payload)
	if err != nil {
		return err
	}
	_, err = wsConn.conn.Write(frame.Bytes())
	if err != nil {
		return err
	}
	bufferPool.Put(frame)
	return nil
}

func (wsConn *WsConn) createFrame(opcode Opcode, payload []byte) (*bytes.Buffer, error) {
	n := len(payload)
	buf := bufferPool.Get(headerFrameLen + n)
	f := &Frame{}
	headerLen, maskingKey := f.CreateHeader(true, opcode, !wsConn.server, n)

	buf.Write(f.Header[:headerLen])
	if !wsConn.server {
		// 客户端需要 掩码
		UnMaskPayload(payload, maskingKey)
	}
	buf.Write(payload)
	return buf, nil
}
