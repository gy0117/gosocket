package main

import (
	"encoding/json"
	"github.com/gy/gosocket"
	"github.com/gy/gosocket/internal/cmap"
	"log"
)

type WebSocketHandler struct {
	sessionMap *cmap.ConcurrentMap[string, *gosocket.WsConn] // 当前服务器管理连接的
}

func (w *WebSocketHandler) OnStart(conn *gosocket.WsConn) {
	sm := conn.GetSessionMap()
	name, _ := sm.Get("name")
	log.Println("WebSocketHandler ---> OnStart, name: ", name)
	c, ok := w.sessionMap.Get(name.(string))
	if ok {
		log.Printf("connection %v has existed, name: %s\n", c, name)
		return
	}
	w.sessionMap.Put(name.(string), conn)
}

func (w *WebSocketHandler) OnPing(conn *gosocket.WsConn, payload []byte) {
	log.Println("WebSocketHandler ---> OnPing")
	err := conn.WriteString("pong")
	if err != nil {
		log.Println("WebSocketHandler ---> OnPing, error: ", err)
	}
}

func (w *WebSocketHandler) OnPong(conn *gosocket.WsConn, payload []byte) {
	log.Println("WebSocketHandler ---> OnPong")
}

func (w *WebSocketHandler) OnMessage(conn *gosocket.WsConn, msg *gosocket.Message) {
	sm := conn.GetSessionMap()
	from, _ := sm.Get("name")

	log.Println("WebSocketHandler ---> OnMessage, from_name: ", from)
	// 在 Chrome 浏览器中使用 WebSocket 时，确实没有暴露给开发者直接发送 Ping 帧的 API。
	b := msg.Content.Bytes()
	if len(b) == 4 && string(b) == "ping" {
		w.OnPing(conn, nil)
		return
	}

	var input = &Input{}
	_ = json.Unmarshal(msg.Content.Bytes(), input)
	if c, ok := w.sessionMap.Get(input.To); ok {
		_ = c.WriteMessage(gosocket.OpcodeTextFrame, msg.Content.Bytes())
		log.Printf("WebSocketHandler ---> OnMessage, opcode=%d, from=%s, to=%s, msg=%s\n", gosocket.OpcodeTextFrame, from, input.To, string(msg.Content.Bytes()))
	}
}

func (w *WebSocketHandler) OnStop(conn *gosocket.WsConn, err error) {
	log.Println("WebSocketHandler ---> OnStop")
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		sessionMap: cmap.New[string, *gosocket.WsConn](10, 64),
	}
}

type Input struct {
	To   string `json:"to"`
	Text string `json:"text"`
}
