package query

import (
	"github.com/graphql-go/graphql"
	"messenger/model"
	"errors"
	"fmt"
)

var Query = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        model.UserType,
				Description: "Get user by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, errors.New("invalid id")
					}

					user := &model.User{
						Id: int64(id),
					}

					result, err := user.Load()

					if err != nil {
						return nil, err
					}
					result.Password = ""

					return result, err
				},
			},

			"users": &graphql.Field{
				Type:        graphql.NewList(model.UserType),
				Description: "Get user list",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 50,
					},
					"skip": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					auth := model.GetAuth(params)

					fmt.Println("Auth", auth)

					limit := params.Args["limit"].(int)
					skip := params.Args["skip"].(int)

					users, err := model.Users(limit, skip)

					if err != nil {
						return nil, err
					}
					return users, err
				},
			},
			"countUsers": &graphql.Field{
				Type:        graphql.Int,
				Description: "Get user list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					count, err := model.CountUsers()

					if err != nil {
						return nil, err
					}
					return count, err
				},
			},
			"message": &graphql.Field{
				Type:        model.MessageType,
				Description: "Get message by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, errors.New("invalid id")
					}

					message := &model.Message{
						Id: int64(id),
					}

					result, err := message.Load()

					if err != nil {
						return nil, err
					}

					return result, err
				},
			},
			"messages": &graphql.Field{
				Type: graphql.NewList(model.MessageType),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 50,
					},
					"skip": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Description: "Get messages list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					limit := params.Args["limit"].(int)
					skip := params.Args["skip"].(int)

					messages, err := model.Messages(limit, skip)
					if err != nil {
						return nil, err
					}
					return messages, err
				},
			},
			"group": &graphql.Field{
				Type: model.GroupType,
				Args: graphql.FieldConfigArgument{

					"id": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Description: "Get group list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, ok := params.Args["id"].(int)

					if !ok {
						return nil, errors.New("invalid group id")
					}

					var userId int64

					auth := model.GetAuth(params)

					if auth != nil {
						userId = auth.UserId
					}
					result, err := model.Load(int64(id), userId)

					if err != nil {
						return nil, err
					}
					return result, err
				},
			},
			"groups": &graphql.Field{
				Type: graphql.NewList(model.GroupType),
				Args: graphql.FieldConfigArgument{
					"search": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"user_id": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 50,
					},
					"skip": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Description: "Get group list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					search := params.Args["search"].(string)
					limit := params.Args["limit"].(int)
					skip := params.Args["skip"].(int)
					uid, ok := params.Args["user_id"].(int)
					auth := model.GetAuth(params)

					var userId int64
					if auth != nil {
						userId = auth.UserId
					}

					if ok {
						userId = int64(uid)
					}

					fmt.Println("userId", userId)
					result, err := model.Groups(search, userId, limit, skip)
					if err != nil {
						return nil, err
					}
					return result, err
				},
			},
			"unread": &graphql.Field{
				Type: graphql.NewList(model.MessageType),
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 50,
					},
					"skip": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Description: "Get unread message",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					limit := params.Args["limit"].(int)
					skip := params.Args["skip"].(int)
					uid, ok := params.Args["user_id"].(int)
					auth := model.GetAuth(params)

					var userId int64
					if auth != nil {
						userId = auth.UserId
					}

					if ok {
						userId = int64(uid)
					}

					fmt.Println("userId", userId)
					result, err := model.UnreadMessages(userId, limit, skip)
					if err != nil {
						return nil, err
					}
					return result, err
				},
			},
		},
	})
