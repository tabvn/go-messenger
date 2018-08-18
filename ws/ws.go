package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"messenger/db"
	"database/sql"
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
	Message json.RawMessage `json:"message"`
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
	} else {
		c.UserId = 0
	}

}
func (w *Ws) RemoveClient(c *Client) {
	delete(w.Clients, c.Id)
}

func (w *Ws) OnMessage(c *Client, message []byte) {

	var m Message

	if err := json.Unmarshal(message, &m); err != nil {

		return
	}

	switch m.Action {

	case AUTH:

		var auth AuthMessage

		err := json.Unmarshal(m.Message, &auth)

		if err == nil {
			w.AuthClient(c, auth.Token)
		}

		break

	default:
		break
	}

}
