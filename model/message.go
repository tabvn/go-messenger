package model

import (
	"time"
	"github.com/graphql-go/graphql"
	"errors"
	"messenger/db"
	"database/sql"
)

type Message struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	GroupId     int64  `json:"group_id"`
	Body        string `json:"body"`
	Emoji       bool   `json:"emoji"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
	Attachments [] Attachment
	Gifs        [] Gif
}

var MessageType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Message",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"group_id": &graphql.Field{
				Type: graphql.Int,
			},
			"body": &graphql.Field{
				Type: graphql.String,
			},
			"emoji": &graphql.Field{
				Type: graphql.Boolean,
			},
			"created": &graphql.Field{
				Type: graphql.Int,
			},
			"updated": &graphql.Field{
				Type: graphql.Int,
			},
			"attachments": &graphql.Field{
				Type: graphql.NewList(AttachmentType),
			},
			"gifs": &graphql.Field{
				Type: graphql.NewList(GifType),
			},
		},
	},
)

func (m *Message) Create() (*Message, error) {

	query := `INSERT INTO messages (user_id, group_id, body, emoji, created, updated) VALUES (?, ?, ?, ?, ?, ?)`

	currentTime := time.Now()
	unixTime := currentTime.Unix()

	m.Created = unixTime
	m.Updated = unixTime

	insertedId, err := db.DB.Insert(query, m.UserId, m.GroupId, m.Body, m.Emoji, m.Created, m.Updated)

	if err != nil {
		return nil, err
	}

	for index, attachment := range m.Attachments {

		q := `INSERT INTO attachments (user_id, message_id, name, original, type, size, created) VALUES (?, ?, ?, ?, ?, ?, ?)`

		attachment.Created = unixTime

		attachmentId, err := db.DB.Insert(q, m.UserId, insertedId, attachment.Name, attachment.Original, attachment.Type, attachment.Size, attachment.Created)

		if err == nil {
			attachment.Id = attachmentId
			attachment.MessageId = insertedId

			m.Attachments[index] = attachment
		}

	}

	m.Id = insertedId

	return m, err
}

func (m *Message) Update() (*Message, error) {

	currentTime := time.Now()
	m.Updated = currentTime.Unix()

	query := `UPDATE users SET body=?, created=? WHERE id = ?`
	_, err := db.DB.Update(query, m.Body, m.Updated, m.Id)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func scanMessage(s db.RowScanner) (*Message, error) {

	var (
		id      int64
		userId  int64
		groupId int64
		body    sql.NullString
		emoji   bool
		created sql.NullInt64
		updated sql.NullInt64
	)

	if err := s.Scan(&id, &userId, &groupId, &body, &emoji, &created, &updated);
		err != nil {
		return nil, err
	}

	m := &Message{
		Id:      id,
		UserId:  userId,
		GroupId: groupId,
		Emoji:   emoji,
		Created: created.Int64,
		Updated: updated.Int64,
	}
	return m, nil
}

func (u *Message) Load() (*Message, error) {

	row, err := db.DB.Get("users", u.Id)
	if err != nil {
		return nil, err
	}

	message, err := scanMessage(row)

	if message == nil {
		return nil, errors.New("user not found")
	}

	return message, err
}

func (m *Message) Delete() (bool, error) {

	_, err := db.DB.Delete("DELETE FROM messages where id=?", m.Id)

	if err != nil {

		return false, err
	}

	return true, nil
}
