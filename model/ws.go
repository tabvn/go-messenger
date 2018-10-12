package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"messenger/db"
	"database/sql"
	"net/http"
	"github.com/satori/go.uuid"
	"log"
	"messenger/helper"
)

var upgrader = websocket.Upgrader{}

const (
	AUTH         = "auth"
	CALL         = "call"
	CALLEND      = "call_end"
	CALLJOIN     = "call_join"
	CALLEXCHANGE = "call_exchange"
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

	case CALL:

		handleReceiveCallMessage(&m)

		break

	case CALLEND:

		handleReceiveCallEnd(&m)

		break

	case CALLJOIN:

		handleReceiveCallJoin(&m)

		break

	case CALLEXCHANGE:

		handleReceiveCallExchange(&m)

		break

	default:
		break
	}

}

func handleReceiveCallExchange(m *Msg) {

	var c map[string]interface{}

	err := json.Unmarshal(m.Payload, &c)

	if err == nil {

		to, ok := c["to"].(float64)
		if !ok {
			return
		}

		userId := int64(to)

		p := map[string]interface{}{
			"action":  CALLEXCHANGE,
			"payload": m.Payload,
		}
		Instance.SendJson(userId, p)

	}
}

func handleReceiveCallJoin(m *Msg) {

	var c map[string]interface{}

	err := json.Unmarshal(m.Payload, &c)

	if err == nil {
		to := helper.GetIds(c["to"])
		gid, k := c["group_id"].(float64)
		if !k {
			return
		}
		uid, ok := c["user_id"].(float64)
		if !ok {
			return
		}
		fr, ok := c["from"].(float64)
		if !ok {
			return
		}

		from := int64(fr)

		userId := int64(uid)
		groupId := int64(gid)
		if err != nil {
			return
		}

		user, err := GetUser(userId)
		if err != nil {
			return
		}

		payload := map[string]interface{}{
			"action": CALLJOIN,
			"payload": map[string]interface{}{
				"group_id": groupId,
				"user": map[string]interface{}{
					"id":         user.Id,
					"first_name": user.FirstName,
					"last_name":  user.LastName,
					"avatar":     user.Avatar,
					"status":     UserStatus(user.Online, user.CustomStatus),
				},
			},
		}

		// send to caller
		Instance.SendJson(from, payload)

		// send to other users
		for _, id := range to {

			if id != from && id != userId {
				Instance.SendJson(id, payload)
			}

		}

	}
}

func handleReceiveCallEnd(m *Msg) {

	var c map[string]interface{}

	err := json.Unmarshal(m.Payload, &c)

	if err == nil {
		to := helper.GetIds(c["to"])
		gid, k := c["group_id"].(float64)
		if !k {
			return
		}
		fr, ok := c["from"].(float64)
		if !ok {
			return
		}

		from := int64(fr)
		groupId := int64(gid)
		if err != nil {
			return
		}

		for _, userId := range to {

			payload := map[string]interface{}{
				"action": CALLEND,
				"payload": map[string]interface{}{
					"group_id": groupId,
					"caller":   from,
				},
			}

			Instance.SendJson(userId, payload)
		}

	}
}

func handleReceiveCallMessage(m *Msg) {

	var c map[string]interface{}

	err := json.Unmarshal(m.Payload, &c)

	if err == nil {
		to := helper.GetIds(c["to"])
		gid, k := c["group_id"].(float64)
		if !k {
			return
		}
		fr, ok := c["from"].(float64)
		if !ok {
			return
		}

		from := int64(fr)
		groupId := int64(gid)

		fromUser, err := GetUser(from)
		if err != nil {
			return
		}

		for _, userId := range to {

			group, e := LoadGroup(groupId, userId)
			if e != nil {
				continue
			}

			payload := map[string]interface{}{
				"action": "calling",
				"payload": map[string]interface{}{
					"group_id": groupId,
					"group":    group,
					"caller": map[string]interface{}{
						"id":         fromUser.Id,
						"first_name": fromUser.FirstName,
						"last_name":  fromUser.LastName,
						"avatar":     fromUser.Avatar,
						"status":     UserStatus(fromUser.Online, fromUser.CustomStatus),
					},
				},
			}

			Instance.SendJson(userId, payload)
		}

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

			break
		}
	}

}
