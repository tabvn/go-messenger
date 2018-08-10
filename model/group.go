package model

import (
	"github.com/graphql-go/graphql"
	"messenger/db"
	"database/sql"
	"fmt"
)

type Group struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Title    string `json:"title"`
	Avatar   string `json:"avatar"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
	Users    [] *User
	Messages [] *Message
}

var GroupType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Group",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"avatar": &graphql.Field{
				Type: graphql.String,
			},
			"created": &graphql.Field{
				Type: graphql.Int,
			},
			"updated": &graphql.Field{
				Type: graphql.Int,
			},
			"users": &graphql.Field{
				Type: graphql.NewList(UserType),
			},
			"messages": &graphql.Field{
				Type: graphql.NewList(MessageType),
			},
		},
	},
)

func scanGroup(rows *sql.Rows) ([] *Group, error) {

	var (
		id                  int64
		userId              int64
		title               sql.NullString
		avatar              sql.NullString
		created             int64
		updated             int64
		messageId           sql.NullInt64
		messageUserId       sql.NullInt64
		messageGroupId      sql.NullInt64
		messageBody         sql.NullString
		messageEmoji        sql.NullBool
		messageCreated      sql.NullInt64
		messageUpdated      sql.NullInt64
		attachmentId        sql.NullInt64
		attachmentMessageId sql.NullInt64
		attachmentName      sql.NullString
		attachmentOriginal  sql.NullString
		attachmentType      sql.NullString
		attachmentSize      sql.NullInt64

		uUserId       sql.NullInt64
		uFirstName    sql.NullString
		uLastName     sql.NullString
		uAvatar       sql.NullString
		uOnline       sql.NullBool
		uCustomStatus sql.NullString
	)

	var groups []*Group
	var group *Group

	for rows.Next() {

		var message *Message
		var attachment *Attachment
		var user *User

		if err := rows.Scan(&id, &userId, &title, &avatar, &created, &updated, &messageId, &messageUserId, &messageGroupId, &messageBody, &messageEmoji,
			&messageCreated, &messageUpdated, &attachmentId, &attachmentMessageId, &attachmentName, &attachmentOriginal, &attachmentType, &attachmentSize,
			&uUserId, &uFirstName, &uLastName, &uAvatar, &uOnline, &uCustomStatus,
		); err != nil {
			fmt.Println("Scan message error", err)
		}

		fmt.Println("scan:", id, messageGroupId.Int64, messageId.Int64)

		if messageId.Int64 > 0 {
			// has message
			message = &Message{
				Id:      messageId.Int64,
				GroupId: messageGroupId.Int64,
				UserId:  messageUserId.Int64,
				Body:    messageBody.String,
				Emoji:   messageEmoji.Bool,
				Created: messageCreated.Int64,
				Updated: messageUpdated.Int64,
			}
		}

		if attachmentId.Int64 > 0 {
			// so we got attachment
			attachment = &Attachment{
				Id:        attachmentId.Int64,
				MessageId: attachmentMessageId.Int64,
				Name:      attachmentName.String,
				Original:  attachmentOriginal.String,
				Type:      attachmentType.String,
				Size:      int(attachmentSize.Int64),
			}

		}
		if uUserId.Int64 > 0 {
			user = &User{
				Id:           uUserId.Int64,
				FirstName:    uFirstName.String,
				LastName:     uLastName.String,
				Avatar:       uAvatar.String,
				Online:       uOnline.Bool,
				CustomStatus: uCustomStatus.String,
			}
		}

		if group != nil && group.Id == id {

			// has group now find all message and insert to group
			// #1 find attachment

			for _, g := range groups {
				for _, m := range g.Messages {
					if attachment != nil && attachmentMessageId.Int64 == m.Id {
						m.Attachments = append(m.Attachments, attachment)
						break
					}

				}

				// check if user in group
				if user != nil {
					hasUser := false

					for _, u := range g.Users {
						if u.Id == user.Id {
							hasUser = true
						}
					}

					if !hasUser {
						g.Users = append(g.Users, user)
					}
				}

			}

		} else {

			group = &Group{
				Id:      id,
				UserId:  userId,
				Title:   title.String,
				Avatar:  avatar.String,
				Created: created,
				Updated: updated,
			}

			if message != nil && messageGroupId.Int64 == id {

				if attachment != nil {
					message.Attachments = append(message.Attachments, attachment)
				}
				group.Messages = append(group.Messages, message)

			}
			if user != nil {
				group.Users = append(group.Users, user)
			}
			groups = append(groups, group)

		}

	}

	return groups, nil
}

func Groups(search string, userId int64, limit int, skip int) ([]*Group, error) {

	var rows *sql.Rows
	var err error
	var query string

	if search == "" {
		query = `
		SELECT g.id, g.user_id, g.title, g.avatar, g.created, g.updated, message.id, message.user_id, message.group_id, 
		message.body, message.emoji, message.created, message.updated, a.id, a.message_id, a.name, a.original, a.type,
		a.size, u.id, u.first_name, u.last_name, u.avatar, u.online, u.custom_status FROM groups as g 
		INNER JOIN members AS m ON m.group_id = g.id AND m.blocked = 0 AND m.user_id = ? 
		LEFT JOIN users as u ON u.id = m.user_id LEFT JOIN messages as message ON message.group_id = g.id 
		AND message.id = (SELECT MAX(id) FROM messages WHERE group_id = g.id ) 
		LEFT JOIN attachments as a ON a.message_id = message.id ORDER BY message.id ASC LIMIT ? OFFSET ?
	`

		rows, err = db.DB.List(query, userId, limit, skip)
	} else {

		search = "%" + search + "%"
		query = `
		SELECT g.id, g.user_id, g.title, g.avatar, g.created, g.updated, message.id, message.user_id, message.group_id, 
		message.body, message.emoji, message.created, message.updated, a.id, a.message_id, a.name, a.original, a.type,
		a.size, u.id, u.first_name, u.last_name, u.avatar, u.online, u.custom_status FROM groups as g 
		INNER JOIN members AS m ON m.group_id = g.id AND m.blocked = 0 AND m.user_id = ? 
		LEFT JOIN users as u ON u.id = m.user_id LEFT JOIN messages as message ON message.group_id = g.id 
		AND message.id = (SELECT MAX(id) FROM messages WHERE group_id = g.id ) 
		LEFT JOIN attachments as a ON a.message_id = message.id WHERE 
		g.title like ? OR u.first_name LIKE ? OR u.last_name LIKE ? OR MATCH(message.body) AGAINST(?) 
		ORDER BY message.id ASC LIMIT ? OFFSET ?
		`

		rows, err = db.DB.List(query, userId, search, search, search, search, limit, skip)
	}

	result, err := scanGroup(rows)

	if err != nil {
		return nil, err
	}

	return result, nil

}
