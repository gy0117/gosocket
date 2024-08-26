package main

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var html []byte

func main() {
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(html); err != nil {
			log.Printf("open index.html failed, err: %v\n", err)
			return
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 升级成websocket
		log.Println("call ws")
	})

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("listen %d failed, err: %+v\n", 9090, err)
	}
}
