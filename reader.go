package gosocket

import (
	"bytes"
	"fmt"
	"github.com/gy/gosocket/internal"
	"io"
)

func (wsConn *WsConn) readMessage() error {
	payloadLen, err := wsConn.frame.ParseHeader(wsConn.bufReader)
	if err != nil {
		return err
	}
	if payloadLen > wsConn.config.MaxReadPayloadSize {
		return internal.ErrCloseTooLarge
	}

	// 在大多数情况下，WebSocket 帧的 RSV1、RSV2、RSV3 位的值都是 0，除非使用了特定的 WebSocket 扩展功能，这些位才可能被设置为 1。
	// TODO 暂时不考虑压缩扩展
	if wsConn.frame.GetRSV1() || wsConn.frame.GetRSV2() || wsConn.frame.GetRSV3() {
		return internal.ErrCloseProtocol
	}

	if err := wsConn.checkMask(); err != nil {
		return err
	}

	var opcode = wsConn.frame.GetOpcode()
	if !opcode.IsDataFrame() {
		return wsConn.readControlFrame()
	}
	return nil
}

// 读取控制帧
func (wsConn *WsConn) readControlFrame() error {
	// RFC 6455 控制帧不可以分片
	if wsConn.frame.GetFIN() == 0 {
		return internal.ErrCloseProtocol
	}
	payloadLen := wsConn.frame.GetPayloadLen()
	// RFC 6455 控制帧的负载长度不得超过 125 字节
	if payloadLen > maxControlFramePayloadLen {
		return internal.ErrCloseProtocol
	}

	// 控制帧的负载数据长度通常很短，因此可以直接读取并处理，这里不用buffer
	payload := make([]byte, payloadLen)
	if payloadLen > 0 {
		if _, err := io.ReadFull(wsConn.bufReader, payload); err != nil {
			return err
		}
		mask := wsConn.frame.GetMask()
		if mask {
			// 还原payload
			UnMaskPayload(payload, wsConn.frame.GetMaskingKey())
		}
	}
	opcode := wsConn.frame.GetOpcode()
	if opcode == OpcodePingFrame {
		wsConn.eventHandler.OnPing(wsConn, payload)
	} else if opcode == OpcodePongFrame {
		wsConn.eventHandler.OnPong(wsConn, payload)
	} else if opcode == OpcodeConnectionCloseFrame {
		return wsConn.handleClose(bytes.NewBuffer(payload))
	} else {
		return internal.NewXError(internal.ErrCloseProtocol, fmt.Errorf("unsupported opcode %d", opcode))
	}
	return nil
}

// 检查掩码设置是否符合 RFC6455 协议
// TODO 确认下
func (wsConn *WsConn) checkMask() error {
	maskEnable := wsConn.frame.GetMask()
	// 服务器不掩码，即mask位必须为0
	if wsConn.isServer && maskEnable {
		return internal.ErrCloseProtocol
	}
	// 客户端必须掩码，即mask位必须为1
	if !wsConn.isServer && !maskEnable {
		return internal.ErrCloseProtocol
	}
	return nil
}
