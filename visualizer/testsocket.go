package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID(), s.URL(), s.LocalAddr(), s.RemoteAddr(), s.Namespace())
		return nil
	})
	//
	server.OnEvent("/", "chat message", func(s socketio.Conn, msg string) {
		fmt.Println("RECEIVED CHAT MESSAGE", msg)
		s.Emit("chat message", msg)
		fmt.Println("RE EMITTING CHAT MESSAGE", msg)
	})

	//server.OnEvent("/", "chat message", func(s socketio.Conn) {
	//	fmt.Println("RECEIVED CHAT MESSAGE AAAA")
	//})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	//
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})


	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("../visualizer")))
	log.Println("Serving at localhost:8000...")

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.Write([]byte("{\"hello\": \"world\"}"))
	//})
	//
	//// Use default options
	//handler := cors.Default().Handler(mux)

	//fmt.Println("AM HERE")

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		log.Fatal(http.ListenAndServe(":8000", nil))
		wg.Done()
	}()

	fmt.Println("am here")
	time.Sleep(10*time.Second)
	server.BroadcastToNamespace("/", "chat message", "HELLO")
	fmt.Println("am here 2")
	wg.Wait()
}