package model

import (
	"time"
	"github.com/graphql-go/graphql"
	"messenger/db"
	"database/sql"
	"fmt"
	"messenger/helper"
	"errors"
	"strconv"
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
	Read        bool   `json:"read"`
	Gif         string `json:"gif"`
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
			"gif": &graphql.Field{
				Type: graphql.String,
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
		},
	},
)

func UserCanSendMessage(userId int64, groupId int64) (bool) {

	q := `SELECT COUNT(*) FROM members WHERE user_id = ? AND group_id = ? AND blocked = 0`

	row, err := db.DB.FindOne(q, userId, groupId)

	if err != nil {
		return false
	}

	var count int64
	if row.Scan(&count) != nil {
		return false
	}

	if count > 0 {
		return true
	}

	return false
}

func (m *Message) Create() (error) {

	if !UserCanSendMessage(m.UserId, m.GroupId) {
		return errors.New("user is not member of group or blocked")
	}

	query := `INSERT INTO messages (user_id, group_id, body, emoji, created, updated) VALUES (?, ?, ?, ?, ?, ?)`

	currentTime := time.Now()
	unixTime := currentTime.Unix()

	m.Created = unixTime
	m.Updated = unixTime

	insertedId, err := db.DB.Insert(query, m.UserId, m.GroupId, m.Body, m.Emoji, m.Created, m.Updated)

	if err != nil {
		return err
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

	return err
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
		gif                 sql.NullString
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

		if err := rows.Scan(&id, &userId, &groupId, &body, &emoji, &gif, &created, &updated, &attachmentId, &attachmentMessageId, &attachmentName, &attachmentOriginal, &attachmentType, &attachmentSize,
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
				Size:      attachmentSize.Int64,
			}
		}

		if message != nil && message.Id == id {
			// exist so need append attachments
			if attachmentId.Int64 > 0{
				message.Attachments = append(message.Attachments, attachment)
			}


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
		a.id, a.message_id, f.name, f.original, f.type, f.size,
		(SELECT COUNT(DISTINCT id) from read_messages WHERE message_id = m.id AND user_id =?) as isRead
		FROM messages AS m LEFT JOIN attachments as a
		ON m.id = a.message_id 
		LEFT JOIN files as f ON a.file_id = f.id
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
		SELECT m.id, m.user_id, m.group_id, m.body, m.emoji, m.gif, m.created, m.updated, 
		a.id, a.message_id, f.name, f.original, f.type, f.size,
		(SELECT COUNT(DISTINCT id) from read_messages WHERE message_id = m.id AND user_id =?) as isRead
		FROM messages AS m 
		LEFT JOIN attachments as a ON m.id = a.message_id
		LEFT JOIN files as f ON a.file_id = f.id
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
		SELECT m.id, m.user_id, m.group_id, m.body, m.emoji, m.gif, m.created, m.updated, 
		a.id, a.message_id, f.name, f.original, f.type, f.size
		FROM messages AS m
		LEFT JOIN attachments as a 
		ON m.id = a.message_id 
		LEFT JOIN files AS f ON a.file_id = f.id
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

func UserCanDeleteMessage(userId, id int64) (bool) {

	q := `SELECT COUNT(*) as c FROM messages WHERE id =? AND user_id =?`

	row, err := db.DB.FindOne(q, id, userId)

	if err != nil {
		return false
	}

	var count sql.NullInt64

	if row.Scan(&count) != nil {
		return false
	}

	if count.Int64 > 0 {
		return true
	}

	return false
}

func MarkMessageAsRead(id int64, userId int64) (error) {

	q := `INSERT INTO read_messages (message_id, user_id, created) VALUES (?, ?, ?)`

	fmt.Println("begin insert read")
	_, err := db.DB.Insert(q, id, userId, helper.GetUnixTimestamp())

	return err
}

func CreateMessage(groupId int64, userId int64, body string, emoji bool, gif string, attachments [] int64) (*Message, error) {

	unixTime := helper.GetUnixTimestamp()
	messageId, err := db.DB.Insert(`INSERT INTO messages (group_id, user_id, body, emoji, gif, created, updated) VALUES (?,?,?,?,?,?,?)`,
		groupId, userId, body, gif, emoji, unixTime, unixTime)

	fmt.Println("got message id created", messageId, err)
	if err != nil {
		return nil, err
	}

	if messageId > 0 {
		// let do insert attachments

		message := &Message{
			Id:      messageId,
			UserId:  userId,
			Body:    body,
			Emoji:   emoji,
			GroupId: groupId,
			Gif:     gif,
			Created: unixTime,
		}

		if len(attachments) > 0 {

			inArrString := ""

			for _, u := range attachments {

				if inArrString == "" {
					inArrString += strconv.Itoa(int(u))
				} else {
					inArrString += ", " + strconv.Itoa(int(u))
				}

			}

			q := fmt.Sprintf("SELECT * FROM files WHERE id IN (%s) AND user_id=?", inArrString)

			rows, err := db.DB.List(q, userId)
			if err != nil {
				return nil, err
			}

			files, err := scanFiles(rows)

			if err != nil {
				return nil, err
			}

			// let create attachments
			value := ""

			for index, file := range files {
				if index == 0 {

					value += fmt.Sprintf("(%d, %d)", messageId, file.Id)
				} else {
					value += fmt.Sprintf(", (%d, %d)", messageId, file.Id)
				}

			}
			if len(files) > 0 {

				for _, file := range files {

					attachmentId, err := db.DB.Insert("INSERT INTO attachments (message_id, file_id) VALUES (?,?)", messageId, file.Id)
					if err == nil && attachmentId > 0 {
						attachment := &Attachment{
							Id:        attachmentId,
							MessageId: messageId,
							Name:      file.Name,
							Original:  file.Original,
							Size:      file.Size,
							Type:      file.Type,
						}
						message.Attachments = append(message.Attachments, attachment)
					}
				}

			}
		}

		return message, nil
	}

	return nil, errors.New("unknown error")
}

func CreateConversation(authorId int64, userIds []int64, messageBody string, messageGif string, messageEmoji bool, attachments [] int64) (*Group, error) {

	gid, err := FindOrCreateGroup(authorId, userIds, "", "")

	if err != nil {
		return nil, err
	}
	if gid > 0 {
		// got group let create message
		_, err := CreateMessage(gid, authorId, messageBody, messageEmoji, messageGif, attachments)

		if err != nil {
			return nil, err
		}

		group, err := LoadGroup(gid, authorId)

		return group, err
	}

	return nil, errors.New("unknown error")
}
