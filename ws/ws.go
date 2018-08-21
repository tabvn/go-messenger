package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"messenger/db"
	"database/sql"
	"messenger/model"
	"fmt"
)

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

type Message struct {
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
		model.UpdateUserStatus(userId.Int64, true)

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

		model.UpdateUserStatus(c.UserId, online)

	}

}

func (w *Ws) OnMessage(c *Client, message []byte) {

	var m Message

	if err := json.Unmarshal(message, &m); err != nil {

		return
	}

	fmt.Println("receive client message", m)
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
