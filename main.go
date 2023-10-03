package main

import (
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func wsHandleText(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsHandleText")
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "closing ws connection")

	// MessageType represents the type of a WebSocket message.
	// See https://tools.ietf.org/html/rfc6455#section-5.6
	msgType, msg, err := c.Read(r.Context())
	if err != nil {
		log.Fatal("error to read message")
		return
	}
	log.Printf("msg type is: %v", msgType)
	log.Printf("the msg is: %s", string(msg))
}

func wsHandleJson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsHandleJson")
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	type wsJson struct {
		Msg      string
		Protocol string
		Status   string
	}
	var v wsJson
	for {
		err = wsjson.Read(r.Context(), c, &v)
		if err != nil {
			break
		}
		log.Printf("received: %v", v)

		if v.Status == "close" {
			c.Write(r.Context(), websocket.MessageText, []byte("bye bye"))
			c.Close(websocket.StatusProtocolError, "close connection")
			break
		} else {
			c.Write(r.Context(), websocket.MessageText, []byte("we still connected"))
		}
	}

}

func main() {
	fmt.Println("Initialize WS-server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	http.HandleFunc("/ws-text", wsHandleText)
	http.HandleFunc("/ws-json", wsHandleJson)
	http.ListenAndServe(":8080", nil)
}
