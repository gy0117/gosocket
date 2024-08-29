package main

import (
	"bufio"
	"fmt"
	"github.com/gy/gosocket"
	"log"
	"net/http"
)

var _ gosocket.EventHandler = (*WsEvent)(nil)

var wsEvent = &WsEvent{}

type WsEvent struct {
}

func (w *WsEvent) OnStart(conn *gosocket.WsConn) {
	log.Println("wsEvent ---> OnStart")
}

func (w *WsEvent) OnPing(conn *gosocket.WsConn, payload []byte) {
	log.Println("wsEvent ---> OnPing")
}

func (w *WsEvent) OnPong(conn *gosocket.WsConn, payload []byte) {
	log.Println("wsEvent ---> OnPong")
}

func (w *WsEvent) OnMessage(conn *gosocket.WsConn, msg *gosocket.Message) {
	log.Println("wsEvent ---> OnMessage")
}

func (w *WsEvent) OnStop(conn *gosocket.WsConn, err error) {
	log.Println("wsEvent ---> OnStop")
}

var upgrade = gosocket.NewUpgrade(wsEvent, &gosocket.ServerOptions{
	BufReaderPool: gosocket.NewPool(func() *bufio.Reader {
		return bufio.NewReaderSize(nil, 1024)
	}),
})

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrade.Upgrade(w, r)
		if err != nil {
			log.Println("Upgrade error: ", err)
			return
		}
		go conn.ReadLoop()
		for i := 0; i < 1000000; i++ {

			err := conn.WriteMessage(gosocket.OpcodeTextFrame, []byte(fmt.Sprintf("abc %d", i)))
			if err != nil {
				log.Println("writer error: ", err)
				break
			}
		}
	})
	_ = http.ListenAndServe("127.0.0.1:7777", nil)
}
