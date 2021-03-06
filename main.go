package main

import (
	"fmt"
	"gim/pkg/websocket"
	"net/http"
)

func ServeWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = r.ParseForm() // 127.0.0.1：8000?user=xxxxx
	if len(r.Form["username"]) > 0 {
		fmt.Println(r.Host, " ", r.Form["username"][0])
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		_, err = fmt.Fprintf(w, "%+v\n", err)
		panic(err)
	}

	client := &websocket.Client{
		Nickname: r.Form["username"][0],
		Conn:     conn,
		Pool:     pool,
	}

	pool.Register <- client
	client.Read()
}

func main() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(pool, w, r)
	})
	_ = http.ListenAndServe(":8000", nil)
}
