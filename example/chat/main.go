package main

import (
	_ "embed"
	"errors"
	"github.com/gy/gosocket"
	"log"
	"net/http"
)

//go:embed index.html
var html []byte

func main() {

	handler := NewWebSocketHandler()
	upgrade := gosocket.NewUpgrade(handler, &gosocket.ServerOptions{
		PreSessionHandle: func(r *http.Request, sm gosocket.SessionManager) error {
			name := r.URL.Query().Get("name")
			if len(name) == 0 {
				return errors.New("name is empty")
			}
			sm.Put("name", name)
			secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
			sm.Put("sec-websocket-key", secWebSocketKey)
			return nil
		},
	})

	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(html); err != nil {
			log.Printf("open index.html failed, err: %v\n", err)
			return
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 升级成websocket
		log.Println("call ws")

		conn, err := upgrade.Upgrade(w, r)
		if err != nil {
			log.Printf("call ws failed: %+v\n", err)
			return
		}
		conn.ReadLoop()
	})

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("listen %d failed, err: %+v\n", 8888, err)
	}
}
