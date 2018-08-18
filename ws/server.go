package ws

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	id := uuid.Must(uuid.NewV4()).String()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	// add new client
	client := &Client{
		Id:   id,
		Conn: c,
	}

	Instance.AddClient(client)

	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			// handle remove client and subscriptions
			Instance.RemoveClient(client)
			client = nil

			break
		}

		Instance.OnMessage(client, message)

		for _, i := range Instance.Clients {

			log.Println("c", i.UserId)
		}

		if err != nil {
			log.Println("error:", err)
			break
		}
	}

}
