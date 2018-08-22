package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"messenger/db"
	"database/sql"
	"net/http"
	"github.com/satori/go.uuid"
	"log"
)

var upgrader = websocket.Upgrader{}

const (
	AUTH = "auth"
)

type AuthMessage struct {
	Token string `json:"token"`
}

type Client struct {
	Id     string
	Conn   *websocket.Conn
	UserId int64
}

type Msg struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type Ws struct {
	Clients map[string]*Client
}

var Instance = &Ws{
	Clients: map[string]*Client{},
}

func (w *Ws) AddClient(c *Client) {
	w.Clients[c.Id] = c

}

func (w *Ws) AuthClient(c *Client, token string) {

	var userId sql.NullInt64

	q := `SELECT user_id FROM tokens WHERE token =? `
	row, err := db.DB.FindOne(q, token)
	if err != nil {
		return
	}

	row.Scan(&userId)

	if userId.Valid && userId.Int64 > 0 {
		c.UserId = userId.Int64

		UpdateUserStatus(c.UserId, true, "")

	} else {
		c.UserId = 0
	}

}
func (w *Ws) RemoveClient(c *Client) {

	delete(w.Clients, c.Id)

	// update user status
	if c.UserId > 0 {
		// need update user status
		var online = false
		for _, i := range w.Clients {
			if i.UserId == c.UserId {
				online = true
				break
			}
		}

		UpdateUserStatus(c.UserId, online, "")

	}

}

func (w *Ws) OnMessage(c *Client, message []byte) {

	var m Msg

	if err := json.Unmarshal(message, &m); err != nil {

		return
	}

	switch m.Action {

	case AUTH:

		var auth AuthMessage

		err := json.Unmarshal(m.Payload, &auth)

		if err == nil {
			w.AuthClient(c, auth.Token)
		}

		break

	default:
		break
	}

}

func (w *Ws) Send(userId int64, message []byte) {

	for _, client := range w.Clients {

		if client.UserId == userId {
			client.Conn.WriteMessage(1, message)
		}
	}
}

func (w *Ws) SendJson(userId int64, message interface{}) {
	for _, client := range w.Clients {

		if client.UserId == userId {
			client.Conn.WriteJSON(message)
		}
	}
}

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

			// handle remove client and subscriptions
			Instance.RemoveClient(client)
			client = nil

			break
		}

		Instance.OnMessage(client, message)

		if err != nil {
			log.Println("error:", err)
			break
		}
	}

}
