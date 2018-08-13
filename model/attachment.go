package model

import "github.com/graphql-go/graphql"

type Attachment struct {
	Id        int64  `json:"id"`
	MessageId int64  `json:"message_id"`
	Name      string `json:"name"`
	Original  string `json:"original"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
}

var AttachmentType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Attachment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
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
		},
	},
)
