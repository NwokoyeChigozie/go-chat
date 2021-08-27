package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/socket.io/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		// server.JoinRoom("/socket.io/", "slack", s)
		s.Join("slack")

		// s.Send("slack", "message", args ...interface{})
		return nil
	})

	server.OnEvent("/socket.io/", "salutations", func(s socketio.Conn, msg string) {
		fmt.Println("salutations: " + msg)
		// s.Send("slack", "chat message", msg)
		server.BroadcastToRoom("/socket.io/", "slack", "message", msg)
		// s.Emit("message", msg)
	})

	// server.OnEvent("/socket.io/", "salutations", func(s socketio.Conn, msg string) {
	// 	fmt.Println("salutations: " + msg)

	// 	s.Emit("message", msg)
	// })

	// server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
	// 	fmt.Println("notice:", msg)
	// 	s.Emit("reply", "have "+msg)
	// })

	// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	// 	s.SetContext(msg)
	// 	return "recv " + msg
	// })

	// server.OnEvent("/", "bye", func(s socketio.Conn) string {
	// 	last := s.Context().(string)
	// 	s.Emit("bye", last)
	// 	s.Close()
	// 	return last
	// })

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	// server.OnDisconnect("/", func(s socketio.Conn, reason string) {
	// 	fmt.Println("closed", reason)
	// })

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./views/chat/")))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Println("Serving at localhost:" + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
