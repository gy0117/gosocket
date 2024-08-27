package gosocket

type Message struct {
}

type EventHandler interface {
	// OnStart 建立连接事件
	OnStart(conn *WsConn)
	// OnPing ping
	OnPing(conn *WsConn, payload []byte)
	// OnPong pong
	OnPong(conn *WsConn, payload []byte)
	// OnMessage 发送消息
	OnMessage(conn *WsConn, msg *Message)
	// OnStop 关闭连接
	OnStop(conn *WsConn, err error)
}
