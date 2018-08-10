package model

import "github.com/graphql-go/graphql"

type Gif struct {
	Id        int64  `json:"id"`
	MessageId int64  `json:"message_id"`
	GifId     string `json:"gif_id"`
	Title     string `json:"title"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int    `json:"size"`
	Url       string `json:"url"`
	Preview   string `json:"preview"`
	Created   int64  `json:"created"`
	Updated   int64  `json:"updated"`
}

var GifType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Gif",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},

			"message_id": &graphql.Field{
				Type: graphql.Int,
			},
			"gif_id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"width": &graphql.Field{
				Type: graphql.Int,
			},
			"height": &graphql.Field{
				Type: graphql.Int,
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"preview": &graphql.Field{
				Type: graphql.String,
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
