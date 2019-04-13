package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"messenger/db"
	"messenger/helper"
	"messenger/sanitize"
	"strconv"
	"time"
)

type Message struct {
	Id          int64          `json:"id"`
	UserId      int64          `json:"user_id"`
	GroupId     int64          `json:"group_id"`
	Body        string         `json:"body"`
	Emoji       bool           `json:"emoji"`
	Created     int64          `json:"created"`
	Updated     int64          `json:"updated"`
	Attachments [] *Attachment `json:"attachments"`
	Unread      bool           `json:"unread"`
	Gif         string         `json:"gif"`
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
			"unread": &graphql.Field{
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
		unread              sql.NullBool
	)

	var messages []*Message
	var message *Message

	for rows.Next() {

		if err := rows.Scan(&id, &userId, &groupId, &body, &emoji, &gif, &created, &updated, &attachmentId, &attachmentMessageId, &attachmentName, &attachmentOriginal, &attachmentType, &attachmentSize,
			&unread); err != nil {

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
			if attachmentId.Int64 > 0 {
				message.Attachments = append(message.Attachments, attachment)
			}

		} else {
			message = &Message{
				Id:      id,
				UserId:  userId,
				GroupId: groupId,
				Body:    body.String,
				Emoji:   emoji,
				Gif:     gif.String,
				Created: created,
				Updated: updated,
				Unread:  unread.Bool,
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

	query := `SELECT m.id, m.user_id, m.group_id, m.body, m.emoji, m.gif, m.created, m.updated, 
		a.id, a.message_id, f.name, f.original, f.type, f.size,
		(SELECT COUNT(DISTINCT id) from unreads WHERE message_id = m.id AND user_id = m.user_id) as unread
		FROM messages AS m
		LEFT JOIN attachments as a
		ON m.id = a.message_id 
		LEFT JOIN files as f ON a.file_id = f.id
		WHERE m.id=?`
	rows, err := db.DB.List(query, m.Id)

	messages, err := scanMessage(rows)

	defer rows.Close()

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
		(SELECT COUNT(DISTINCT id) from unreads WHERE message_id = m.id AND user_id = m.user_id) as unread
		FROM messages AS m 
		LEFT JOIN unreads as un ON m.id = un.message_id AND un.user_id =?
		LEFT JOIN attachments as a ON m.id = a.message_id
		LEFT JOIN files as f ON a.file_id = f.id
		INNER JOIN (SELECT mm.id FROM messages as mm WHERE mm.group_id =? ORDER BY mm.id DESC LIMIT ? OFFSET ?) as mj on mj.id = m.id
		WHERE m.user_id NOT IN (SELECT user FROM blocked WHERE author =? AND user = m.user_id)
		AND m.id NOT IN (SELECT message_id FROM deleted WHERE user_id =?  AND message_id = m.id)
	`

	rows, err := db.DB.List(query, userId, groupId, limit, skip, userId, userId)

	messages, err := scanMessage(rows)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func UnreadMessages(userId int64, limit int, skip int) ([] *Message, error) {
	// fake un.id for scan
	query := `
		SELECT m.id, m.user_id, m.group_id, m.body, m.emoji, m.gif, m.created, m.updated, 
		a.id, a.message_id, f.name, f.original, f.type, f.size,
		un.id
		FROM messages AS m
		LEFT JOIN attachments as a 
		ON m.id = a.message_id 
		LEFT JOIN files AS f ON a.file_id = f.id
		INNER JOIN unreads un ON m.id = un.message_id AND un.user_id = ?
		order by m.created DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := db.DB.List(query, userId, limit, skip)

	messages, err := scanMessage(rows)
	defer rows.Close()

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

func UserCanReadMessage(userId, messageId int64) (bool) {

	q := `SELECT COUNT(m.id) as count FROM messages as m INNER JOIN groups as g ON g.id = m.group_id INNER JOIN members as mb ON g.id = mb.group_id AND mb.blocked = 0 AND mb.user_id = ? WHERE m.id = ?`

	row, err := db.DB.FindOne(q, userId, messageId)
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

func UpdateMessage(userId, messageId int64, body string, emoji bool) (bool) {

	updated := helper.GetUnixTimestamp()
	query := `UPDATE messages SET body=?, emoji=?, updated=? WHERE id = ? AND user_id =?`

	body = sanitize.HTML(body)

	rowEffect, err := db.DB.Update(query, body, emoji, updated, messageId, userId)

	if rowEffect > 0 && err == nil {

		defer func() {

			groupId, e := GetMessageGroupId(messageId)
			if groupId > 0 && e == nil {
				ids := GetGroupMemberOnline(userId, groupId)

				payload := map[string]interface{}{
					"action": "message_updated",
					"payload": map[string]interface{}{
						"id":       messageId,
						"user_id":  userId,
						"body":     body,
						"emoji":    emoji,
						"updated":  updated,
						"group_id": groupId,
					},
				}
				for _, id := range ids {
					Instance.SendJson(id, payload)
				}
			}

		}()

		return true
	}

	return false
}

func GetMessageGroupId(messageId int64) (int64, error) {

	findQuery := "SELECT m.group_id FROM messages AS m WHERE m.id =?"

	var scanGroupId sql.NullInt64

	findRow, findErr := db.DB.FindOne(findQuery, messageId)

	if findErr != nil {
		return 0, findErr
	}

	e := findRow.Scan(&scanGroupId)
	if e != nil {
		return 0, e
	}

	groupId := scanGroupId.Int64

	return groupId, nil
}

func DeleteMessage(userId, messageId int64) (bool) {

	// find Message

	groupId, e := GetMessageGroupId(messageId)

	if groupId == 0 || e != nil {
		return false
	}

	// delete the message if is owner
	deleteQuery := `DELETE FROM messages WHERE id=? AND user_id =?`

	numRowDeleted, err := db.DB.Delete(deleteQuery, messageId, userId)

	if err != nil {
		return false
	}

	if numRowDeleted > 0 {
		// this is delete by user so we do need notify message is delete

		defer func() {

			payload := map[string]interface{}{
				"action":  "message_deleted",
				"payload": messageId,
			}

			ids := GetGroupMemberOnline(userId, groupId)

			for _, id := range ids {
				Instance.SendJson(id, payload)
			}

		}()

		return true
	} else {

		isMember := IsMemberOfGroup(userId, groupId)

		if isMember {
			// add to deleted table

			addQuery := "INSERT INTO deleted (user_id, message_id) VALUES (?, ?)"

			r, e := db.DB.Insert(addQuery, userId, messageId)

			if e != nil {
				return false
			}
			if r > 0 {
				return true
			}
		}

	}

	return false
}

func MarkMessageAsRead(id int64, userId int64) (error) {

	q := `DELETE FROM reads WHERE message_id =? AND user_id =?`
	_, err := db.DB.Delete(q, id, userId)

	return err
}

func MarkAsReadByGroup(groupId int64, userId int64) (error) {

	q := `DELETE u.* FROM unreads as u INNER JOIN messages as m ON u.message_id = m.id AND m.group_id =? WHERE u.user_id =?`

	_, err := db.DB.DeleteMany(q, groupId, userId)

	return err
}

func GetGroupMemberOnline(userId, groupId int64) ([] int64) {

	var ids [] int64

	q := `SELECT DISTINCT(m.user_id) FROM members AS m INNER JOIN users AS u ON m.user_id = u.id AND u.online = 1 WHERE m.group_id=? AND m.blocked = 0 AND m.accepted <> 2 AND m.user_id NOT IN (SELECT b.user FROM blocked AS b WHERE b.author =?) AND m.user_id NOT IN (SELECT b.author FROM blocked AS b WHERE b.user =?)`

	rows, err := db.DB.List(q, groupId, userId, userId)
	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var id sql.NullInt64

		if rows.Scan(&id) != nil {
			return nil
		}

		if id.Int64 > 0 {
			ids = append(ids, id.Int64)
		}

	}
	return ids
}
func NotifyMessageToMembers(groupId int64, message Message) {

	payload := map[string]interface{}{
		"action":  "message",
		"payload": message,
	}

	ids := GetGroupMemberOnline(message.UserId, groupId)

	for _, id := range ids {
		Instance.SendJson(id, payload)
	}

}
func CreateMessage(groupId int64, userId int64, body string, emoji bool, gif string, attachments [] int64) (*Message, error) {

	// we need check if members in group is more than 1

	body = sanitize.HTML(body)
	row, err := db.DB.FindOne("SELECT COUNT(*) FROM members WHERE group_id =? AND user_id=?", groupId, userId)
	if err != nil {
		return nil, err
	}
	var count int64

	if row.Scan(&count) != nil {
		return nil, errors.New("can not send message")
	}

	if count < 1 {
		return nil, errors.New("your message could not be sent")
	}

	// we may check and update members accepted

	db.DB.Update("UPDATE members SET accepted = 1 WHERE user_id=? AND group_id=?", userId, groupId)

	unixTime := helper.GetUnixTimestamp()
	messageId, err := db.DB.Insert(`INSERT INTO messages (group_id, user_id, body, emoji, gif, created, updated) VALUES (?,?,?,?,?,?,?)`,
		groupId, userId, body, emoji, gif, unixTime, unixTime)

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
			Updated: unixTime,
			Unread:  true,
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

			defer rows.Close()

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

		defer NotifyMessageToMembers(groupId, *message)

		return message, nil
	}

	return nil, errors.New("unknown error")
}

func CreateConversation(authorId int64, userIds []int64, messageBody string, messageGif string, messageEmoji bool, attachments [] int64, title, avatar string) (*Group, error) {

	gid, err := FindOrCreateGroup(authorId, userIds, title, avatar)

	if err != nil {
		return nil, err
	}
	if gid > 0 {
		// got group let create message
		if len(messageBody) > 0 || len(messageGif) > 0 || len(attachments) > 0 {
			_, err := CreateMessage(gid, authorId, messageBody, messageEmoji, messageGif, attachments)

			if err != nil {
				return nil, err
			}
		}

		group, err := LoadGroup(gid, authorId)

		return group, err
	}

	return nil, errors.New("unknown error")
}
