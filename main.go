package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	clients       map[string]*ClientMetaData = make(map[string]*ClientMetaData)
	globalContext                            = context.Background()
)

type ClientMetaData struct {
	NickName            string          `json:"nickname"`
	Msg                 string          `json:"msg"`
	Target              string          `json:"target"`
	Status              string          `json:"status"`
	WebSocketConnection *websocket.Conn `json:"omitempty"`
}

func (m *ClientMetaData) dmMessage() {
	if len(clients) > 0 {
		for _, v := range clients {
			if (m.Target == v.NickName) && (v.Status == "online") {
				v.WebSocketConnection.Write(globalContext, websocket.MessageText, []byte(m.Msg))
				return
			}
		}
		m.WebSocketConnection.Write(globalContext, websocket.MessageText, []byte("coul'd not find your target"))
	}
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[*] WEB-SCOKET HANDLE STARTED")

	newConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	var newClient ClientMetaData
	newClient.NickName = r.URL.Query().Get("nickname")
	newClient.WebSocketConnection = newConn
	newClient.Status = "online"

	var clientId string = uuid.NewString()
	clients[clientId] = &newClient

	log.Printf("[*] new client add: %v", newClient.NickName)
	log.Printf("[*] total clients: %d", len(clients))

	// lifecycle for newConn
	for {
		err = wsjson.Read(globalContext, newConn, &newClient)
		if err != nil {
			break
		}

		if newClient.Status == "offline" {
			newConn.Write(globalContext, websocket.MessageText, []byte("bye bye"))
			newConn.Close(websocket.StatusProtocolError, "close connection")
			delete(clients, clientId)
			break
		} else {
			newClient.dmMessage()
		}
	}
	log.Println("[*] close handler for id:" + clientId)
}

// a schedulle job job to send every connected user
func showAllClientsConnected() {
	for {
		if len(clients) > 0 {
			time.Sleep(40 * time.Second)
			for _, v := range clients {
				jsonRes, err := json.Marshal(clients)
				if err != nil {
					v.WebSocketConnection.Write(globalContext, websocket.MessageText, []byte("could not send all clients connected"))
				}
				v.WebSocketConnection.Write(globalContext, websocket.MessageText, jsonRes)
			}
		}
	}
}

// func wsHandleText(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("wsHandleText")
// 	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
// 		InsecureSkipVerify: true,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer c.Close(websocket.StatusInternalError, "closing ws connection")

// 	// MessageType represents the type of a WebSocket message.
// 	// See https://tools.ietf.org/html/rfc6455#section-5.6
// 	msgType, msg, err := c.Read(r.Context())
// 	if err != nil {
// 		log.Fatal("error to read message")
// 		return
// 	}
// 	log.Printf("msg type is: %v", msgType)
// 	log.Printf("the msg is: %s", string(msg))
// }

func main() {
	log.Println("[*] initialize ws-server")

	go showAllClientsConnected()

	// http.HandleFunc("/ws-text", wsHandleText)
	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		var res []*ClientMetaData
		for _, v := range clients {
			res = append(res, v)
		}
		log.Println("[*] client user-agent" + r.UserAgent())
		json.NewEncoder(w).Encode(res)
	})
	http.HandleFunc("/ws-json", wsHandle)
	http.ListenAndServe(":8080", nil)
}
