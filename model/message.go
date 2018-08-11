package model

import (
	"time"
	"github.com/graphql-go/graphql"
	"messenger/db"
	"database/sql"
	"fmt"
)

type Message struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	GroupId     int64  `json:"group_id"`
	Body        string `json:"body"`
	Emoji       bool   `json:"emoji"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
	Attachments [] *Attachment
	Gifs        [] *Gif
	Read        bool   `json:"read"`
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
			"read": &graphql.Field{
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

		q := `INSERT INTO attachments (message_id, name, original, type, size) VALUES (?, ?, ?, ?, ?, ?, ?)`

		attachmentId, err := db.DB.Insert(q, insertedId, attachment.Name, attachment.Original, attachment.Type, attachment.Size)

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

	query := `UPDATE messages SET body=?, emoji=?, updated=? WHERE id = ?`

	_, err := db.DB.Update(query, m.Body, m.Emoji, m.Updated, m.Id)

	if err != nil {

		return nil, err
	}

	return m, nil
}

func scanMessage(rows *sql.Rows) ([] *Message, error) {

	var (
		id                  int64
		userId              int64
		groupId             int64
		body                sql.NullString
		emoji               bool
		created             int64
		updated             int64
		attachmentId        sql.NullInt64
		attachmentMessageId sql.NullInt64
		attachmentName      sql.NullString
		attachmentOriginal  sql.NullString
		attachmentType      sql.NullString
		attachmentSize      sql.NullInt64
		read                sql.NullBool
	)

	var messages []*Message
	var message *Message

	for rows.Next() {

		if err := rows.Scan(&id, &userId, &groupId, &body, &emoji, &created, &updated, &attachmentId, &attachmentMessageId, &attachmentName, &attachmentOriginal, &attachmentType, &attachmentSize,
			&read); err != nil {
			fmt.Println("Scan message error", err)
		}

		var attachment *Attachment

		if attachmentId.Int64 != 0 {
			attachment = &Attachment{
				Id:        attachmentId.Int64,
				MessageId: attachmentMessageId.Int64,
				Name:      attachmentName.String,
				Original:  attachmentOriginal.String,
				Type:      attachmentType.String,
				Size:      int(attachmentSize.Int64),
			}
		}

		if message != nil && message.Id == id {
			// exist so need append attachments
			message.Attachments = append(message.Attachments, attachment)

		} else {
			message = &Message{
				Id:      id,
				UserId:  userId,
				GroupId: groupId,
				Body:    body.String,
				Emoji:   emoji,
				Created: created,
				Updated: updated,
				Read:    read.Bool,
			}
			if attachment != nil {
				message.Attachments = append(message.Attachments, attachment)
			}
			messages = append(messages, message)

		}

	}

	return messages, nil
}

func (m *Message) Load() (*Message, error) {

	query := `
		SELECT m.*, 
		a.id, a.message_id, a.name, a.original, a.type, a.size,
		(SELECT COUNT(DISTINCT id) from read_messages WHERE message_id = m.id AND user_id =?) as isRead
		FROM messages AS m LEFT JOIN attachments as a
		ON m.id = a.message_id 
		WHERE m.id=?
		order by m.created DESC, a.id DESC
	`
	rows, err := db.DB.List(query, m.UserId, m.Id)

	messages, err := scanMessage(rows)

	if err != nil {
		return nil, err
	}

	if len(messages) > 0 {
		return messages[0], nil
	}

	return nil, err
}

func Messages(groupId int64, userId int64, limit int, skip int) ([] *Message, error) {

	query := `
		SELECT m.id, m.user_id, m.group_id, m.body, m.emoji, m.created, m.updated, 
		a.id, a.message_id, a.name, a.original, a.type, a.size,
		(SELECT COUNT(DISTINCT id) from read_messages WHERE message_id = m.id AND user_id =?) as isRead
		FROM messages AS m LEFT JOIN attachments as a 
		ON m.id = a.message_id 
		INNER JOIN (SELECT mm.id FROM messages as mm WHERE mm.group_id =? ORDER BY mm.id DESC LIMIT ? OFFSET ?) as mj on mj.id = m.id
	`

	rows, err := db.DB.List(query, userId, groupId, limit, skip)

	messages, err := scanMessage(rows)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func UnreadMessages(userId int64, limit int, skip int) ([] *Message, error) {
	query := `
		SELECT m.id, m.user_id, m.group_id, m.body, m.emoji, m.created, m.updated, 
		a.id, a.message_id, a.name, a.original, a.type, a.size
		FROM messages AS m
		LEFT JOIN attachments as a 
		ON m.id = a.message_id 
		WHERE m.id NOT IN (SELECT message_id FROM read_messages WHERE user_id = ?)
		order by m.created DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := db.DB.List(query, userId, limit, skip)

	messages, err := scanMessage(rows)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *Message) Delete() (bool, error) {

	_, err := db.DB.Delete("DELETE FROM messages where id=?", m.Id)

	if err != nil {

		return false, err
	}

	return true, nil
}
