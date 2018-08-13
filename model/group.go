package model

import (
	"github.com/graphql-go/graphql"
	"messenger/db"
	"database/sql"
	"fmt"
	"errors"
	"messenger/helper"
	"strconv"
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
	Unread   int64  `json:"unread"`
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
			"unread": &graphql.Field{
				Type: graphql.Int,
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
		unread              int64
		messageId           sql.NullInt64
		messageUserId       sql.NullInt64
		messageGroupId      sql.NullInt64
		messageBody         sql.NullString
		messageEmoji        sql.NullBool
		messageIsRead       sql.NullBool
		messageCreated      sql.NullInt64
		messageUpdated      sql.NullInt64
		attachmentId        sql.NullInt64
		attachmentMessageId sql.NullInt64
		attachmentName      sql.NullString
		attachmentOriginal  sql.NullString
		attachmentType      sql.NullString
		attachmentSize      sql.NullInt64

		uUserId       sql.NullInt64
		uid           sql.NullInt64
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
			&uUserId, &uid, &uFirstName, &uLastName, &uAvatar, &uOnline, &uCustomStatus, &unread, &messageIsRead,
		);
			err != nil {
			fmt.Println("Scan message error", err)
		}

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
				Read:    messageIsRead.Bool,
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
				Size:      attachmentSize.Int64,
			}

		}
		if uUserId.Int64 > 0 {
			user = &User{
				Id:           uUserId.Int64,
				Uid:          uid.Int64,
				FirstName:    uFirstName.String,
				LastName:     uLastName.String,
				Avatar:       uAvatar.String,
				Online:       uOnline.Bool,
				CustomStatus: uCustomStatus.String,
				Status:       UserStatus(uOnline.Bool, uCustomStatus.String),
			}
		}

		if group != nil && group.Id == id {

			// has group now find all message and insert to group
			// #1 find attachment

			for _, g := range groups {

				for _, m := range g.Messages {

					hasAttachment := false

					for _, a := range m.Attachments {

						if a.Id == attachmentId.Int64 {
							hasAttachment = true
						}
					}

					if attachment != nil && attachmentMessageId.Int64 == m.Id && !hasAttachment {
						m.Attachments = append(m.Attachments, attachment)
					}

				}

				// check if user in group
				if uUserId.Valid {
					hasUser := false

					for _, u := range g.Users {

						if u.Id == uUserId.Int64 {
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
				Unread:  unread,
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

func LoadGroup(id int64, userId int64) (*Group, error) {

	var rows *sql.Rows
	query := `
		SELECT g.id, g.user_id, g.title, g.avatar, g.created, g.updated, message.id, message.user_id, message.group_id, 
		message.body, message.emoji, message.created, message.updated, a.id, a.message_id, f.name, f.original, f.type,
		f.size, u.id, u.uid, u.first_name, u.last_name, u.avatar, u.online, u.custom_status,
		(SELECT COUNT(DISTINCT cm.id) 
        FROM messages cm WHERE cm.group_id = g.id AND cm.id NOT IN (SELECT message_id FROM read_messages WHERE message_id = cm.id  AND user_id =? )
       ) as unread,
		r.id as mread
		FROM groups as g 
		INNER JOIN members AS m ON m.group_id = g.id
		LEFT JOIN users as u ON u.id = m.user_id LEFT JOIN messages as message ON message.group_id = g.id 
		AND message.id = (SELECT MAX(id) FROM messages WHERE group_id = g.id ) 
		LEFT JOIN read_messages as r ON r.message_id = message.id AND r.user_id =?
		LEFT JOIN attachments as a ON a.message_id = message.id
		LEFT JOIN files as f ON a.file_id = f.id
		WHERE g.id = ?
		
	`
	rows, err := db.DB.List(query, userId, userId, id)

	result, err := scanGroup(rows)

	if len(result) < 1 {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}

	return result[0], nil

}

func Groups(search string, userId int64, limit int, skip int) ([]*Group, error) {

	var rows *sql.Rows
	var err error
	var query string

	if search == "" {
		query = `
		SELECT g.id, g.user_id, g.title, g.avatar, g.created, g.updated, message.id, message.user_id, message.group_id, 
		message.body, message.emoji, message.created, message.updated, a.id, a.message_id, a.name, a.original, a.type,
		a.size, u.id, u.uid, u.first_name, u.last_name, u.avatar, u.online, u.custom_status,
		(SELECT COUNT(DISTINCT cm.id) 
        FROM messages cm WHERE cm.group_id = g.id AND cm.id NOT IN (SELECT message_id FROM read_messages WHERE message_id = cm.id  AND user_id =? )
       ) as unread,
       	r.id as mread
		FROM groups as g 
		INNER JOIN members AS m ON m.group_id = g.id
		LEFT JOIN users as u ON u.id = m.user_id 
		LEFT JOIN messages as message ON message.group_id = g.id 
		AND message.id = (SELECT MAX(id) FROM messages WHERE group_id = g.id ) 
		LEFT JOIN read_messages as r ON r.message_id = message.id AND r.user_id =?
		LEFT JOIN attachments as a ON a.message_id = message.id INNER JOIN (SELECT gr.id FROM groups as gr INNER JOIN members as mb ON gr.id = mb.group_id AND mb.blocked = 0 
		AND mb.user_id =? INNER JOIN messages as msg ON msg.group_id = gr.id GROUP BY gr.id ORDER BY msg.id DESC LIMIT ? OFFSET ?) as grj ON grj.id = g.id
	`

		rows, err = db.DB.List(query, userId, userId, userId, limit, skip)
	} else {

		searchLike := "%" + search + "%"

		query = `
		SELECT g.id, g.user_id, g.title, g.avatar, g.created, g.updated, message.id, message.user_id, message.group_id, 
		message.body, message.emoji, message.created, message.updated, a.id, a.message_id, a.name, a.original, a.type,
		a.size, u.id, u.uid, u.first_name, u.last_name, u.avatar, u.online, u.custom_status,
		(SELECT COUNT(DISTINCT cm.id) 
        FROM messages cm WHERE cm.group_id = g.id AND cm.id NOT IN (SELECT message_id FROM read_messages WHERE message_id = cm.id  AND user_id =? )
       ) as unread,
       	r.id as mread
		FROM groups as g 
		INNER JOIN members AS m ON m.group_id = g.id
		LEFT JOIN users as u ON u.id = m.user_id 
		LEFT JOIN messages as message ON message.group_id = g.id 
		AND message.id = (SELECT MAX(id) FROM messages WHERE group_id = g.id ) 
		LEFT JOIN read_messages as r ON r.message_id = message.id AND r.user_id =?
		LEFT JOIN attachments as a ON a.message_id = message.id INNER JOIN (SELECT gr.id FROM groups as gr INNER JOIN members as mb ON gr.id = mb.group_id AND mb.blocked = 0 
		AND mb.user_id =? INNER JOIN messages as msg ON msg.group_id = gr.id GROUP BY gr.id ORDER BY msg.id DESC LIMIT ? OFFSET ?) as grj ON grj.id = g.id
		WHERE g.title like ? OR u.first_name LIKE ? OR u.last_name LIKE ? OR MATCH(message.body) AGAINST(?)
		`

		rows, err = db.DB.List(query, userId, userId, userId, limit, skip, searchLike, searchLike, searchLike, search)
	}

	result, err := scanGroup(rows)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func CanJoinGroup(authorId, userId, groupId int64) (bool) {

	// we are not allow self join

	if authorId == userId {
		return false
	}
	query := ` SELECT COUNT(m.user_id) as count FROM groups as g  INNER JOIN members as m ON m.group_id = g.id WHERE g.id = ? AND m.user_id = ?
				
	`

	row, err := db.DB.FindOne(query, groupId, authorId)

	if err != nil {
		return false
	}

	var count int

	if row.Scan(&count) != nil {
		return false
	}

	if count > 0 {

		return true
	}

	return false
}

func JoinGroup(userId, groupId int64) (bool) {

	q := `INSERT INTO members (user_id, group_id, blocked, created) VALUES (?, ?, ?, ?)`

	insertedId, err := db.DB.Insert(q, userId, groupId, 0, helper.GetUnixTimestamp())

	fmt.Print("join", insertedId, err)
	if err != nil {
		return false
	}

	if insertedId > 0 {
		return true
	}
	return false
}

func LeftGroup(userId, groupId int64) (int64, error) {

	q := `DELETE FROM members WHERE user_id = ? AND group_id =?`

	result, err := db.DB.Delete(q, userId, groupId)
	return result, err
}

func FindOrCreateGroup(authorId int64, userIds [] int64, title, avatar string) (int64, error) {

	// find group with all members

	inArrString := ""

	for _, u := range userIds {

		if inArrString == "" {
			inArrString += strconv.Itoa(int(u))
		} else {
			inArrString += ", " + strconv.Itoa(int(u))
		}

	}

	inArrString += ""
	total := len(userIds)
	findQuery := fmt.Sprintf(`select g.id,
	(SELECT COUNT(DISTINCT mb.id) 
	FROM members mb WHERE mb.group_id = g.id) as total
	from members as m
	INNER JOIN groups as g ON m.group_id = g.id
	where m.user_id in (%s)
	group by m.group_id
	having count(distinct m.user_id) = ? AND total = ?`, inArrString)
	row, err := db.DB.FindOne(findQuery, total, total)

	if err != nil {
		return 0, err
	}

	var (
		scanGroupId sql.NullInt64
		scanTotal   sql.NullInt64
	)

	scanErr := row.Scan(&scanGroupId, &scanTotal)

	if scanErr != nil && scanErr != sql.ErrNoRows {
		return 0, scanErr
	}

	groupId := scanGroupId.Int64

	if groupId > 0 {
		// group is exist so load group and return
		return groupId, nil

	} else {

		unixTime := helper.GetUnixTimestamp()

		createGroupQuery := `INSERT INTO groups (userId, title, avatar, created, updated) VALUES (?,?,?,?,?)`

		gid, createErr := db.DB.Insert(createGroupQuery, authorId, title, avatar, unixTime, unixTime)

		if createErr != nil {
			return 0, createErr
		}

		// create members

		values := ""

		for i := 0; i < len(userIds); i++ {
			uid := userIds[i]

			if i == 0 {
				str := fmt.Sprintf("(%d, %d, %d, %d)", uid, groupId, 0, unixTime)
				values += str
			} else {
				values += fmt.Sprintf(", (%d, %d, %d, %d)", uid, groupId, 0, unixTime)
			}

		}

		createMemberQuery := `INSERT INTO members (userId, groupId, blocked, created) values ` + values

		numRows, createMemberErr := db.DB.InsertMany(createMemberQuery)
		if createMemberErr != nil {
			db.DB.Delete(`DELETE groups WHERE id =? `, gid)
			return 0, createMemberErr
		}
		if numRows == int64(len(userIds)) {
			return gid, nil
		} else {
			// delete group
			db.DB.Delete(`DELETE groups WHERE id =? `, gid)
		}

	}

	return 0, nil

}
