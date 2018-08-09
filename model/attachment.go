package model

import "github.com/graphql-go/graphql"

type Attachment struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	MessageId int64  `json:"message_id"`
	Name      string `json:"name"`
	Original  string `json:"original"`
	Type      string `json:"type"`
	Size      int    `json:"size"`
	Created   int64  `json:"created"`
	Updated   int64  `json:"updated"`
}

var AttachmentType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Attachment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"message_id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"original": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"created": &graphql.Field{
				Type: graphql.Int,
			},
			"updated": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
