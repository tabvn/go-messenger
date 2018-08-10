package model

import (
	"github.com/graphql-go/graphql"
	"messenger/db"
	"database/sql"
	"fmt"
)

type Group struct {
	Id      int64  `json:"id"`
	UserId  int64  `json:"user_id"`
	Title   string `json:"title"`
	Avatar  string `json:"avatar"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
	Users   [] User
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
		},
	},
)

func scanGroup(rows *sql.Rows) ([] *Group, error) {

	var (
		id                 int64
		userId             int64
		title              sql.NullString
		avatar             sql.NullString
		created            int64
		updated            int64
		messageId          int64
		messageBody        sql.NullString
		messageEmoji       bool
		messageCreated     int64
		messageUpdated     int64
		attachmentId       int64
		attachmentName     sql.NullString
		attachmentOriginal sql.NullString
		attachmentType     sql.NullString
		attachmentSize     int
	)

	var groups []*Group
	var group *Group

	for rows.Next() {

		if err := rows.Scan(&id, &userId, &title, &avatar, &created, &updated, &messageId, &messageBody, &messageEmoji,
			&messageCreated, &messageUpdated, &attachmentId, &attachmentName, &attachmentOriginal, &attachmentType, &attachmentSize); err != nil {
			fmt.Println("Scan message error", err)
		}

		if group != nil && group.Id == id {
			// exist so need append attachments
			//message.Attachments = append(message.Attachments, *attachment)

		} else {
			group = &Group{
				Id:      id,
				UserId:  userId,
				Title:   title.String,
				Avatar:  avatar.String,
				Created: created,
				Updated: updated,
			}

			groups = append(groups, group)

		}

	}

	return groups, nil
}

func Groups(userId int64, limit int, skip int) ([]*Group, error) {

	query := `
		SELECT g.id, g.user_id, g.title, g.avatar, g.created, g.updated, message.id, message.body, message.emoji, message.created, message.updated, a.id, a.name, a.original, a.type, a.size FROM groups as g INNER JOIN members AS m ON m.group_id = g.id AND m.user_id = ?  LEFT JOIN messages as message ON message.group_id = g.id AND message.id = (SELECT MAX(id) FROM messages WHERE group_id = g.id ) LEFT JOIN attachments as a ON a.message_id = message.id ORDER BY message.id ASC LIMIT ? OFFSET ?
	`

	rows, err := db.DB.List(query, userId, limit, skip)

	result, err := scanGroup(rows)

	if err != nil {
		return nil, err
	}

	return result, nil

}
