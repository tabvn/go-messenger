package mutation

import (
	"github.com/graphql-go/graphql"
	"messenger/model"
	"errors"
	"fmt"
)

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{

		"createUser": &graphql.Field{
			Type:        model.UserType,
			Description: "Create new user",
			Args: graphql.FieldConfigArgument{
				"first_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"last_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				user := &model.User{
					FirstName: params.Args["first_name"].(string),
					LastName:  params.Args["last_name"].(string),
					Email:     params.Args["email"].(string),
					Password:  params.Args["password"].(string),
				}

				result, err := user.Create()

				if err != nil {
					return nil, err
				}

				result.Password = ""

				return result, err

			},
		},

		"updateUser": &graphql.Field{
			Type:        model.UserType,
			Description: "Update user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"first_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"last_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				id, ok := params.Args["id"].(int)
				if !ok {
					return nil, errors.New("invalid id")
				}
				user := model.User{
					Id:        int64(id),
					FirstName: params.Args["first_name"].(string),
					LastName:  params.Args["last_name"].(string),
					Email:     params.Args["email"].(string),
					Password:  params.Args["password"].(string),
				}

				result, err := user.Update()

				if err != nil {
					return nil, err
				}
				return result, err

			},
		},

		"deleteUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				id, ok := params.Args["id"].(int)
				if !ok {
					return nil, errors.New("invalid id")
				}

				user := model.User{
					Id: int64(id),
				}

				result, err := user.Delete()

				if err != nil {
					return nil, errors.New("an error deleting user or not found")
				}
				return result, err

			},
		},

		"login": &graphql.Field{
			Type:        model.LoginType,
			Description: "Login",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				email := params.Args["email"].(string)
				password := params.Args["password"].(string)

				token, user, err := model.LoginUser(email, password)

				if err != nil {
					return nil, err
				}

				r := map[string]interface{}{
					"id":      token.Id,
					"token":   token.Token,
					"created": token.Created,
					"user":    user,
				}

				return r, nil

			},
		},
		"logout": &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Logout",
					Fields: graphql.Fields{
						"success": &graphql.Field{
							Type: graphql.Boolean,
						},
					},
				},
			),
			Description: "Login",
			Args: graphql.FieldConfigArgument{
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				token := params.Args["token"].(string)

				success, err := model.LogoutUser(token)

				if err != nil {
					return nil, errors.New("logout error")
				}
				result := map[string]interface{}{
					"success": success,
				}

				return result, err

			},
		},
		"createMessage": &graphql.Field{
			Type:        model.MessageType,
			Description: "Create new message",
			Args: graphql.FieldConfigArgument{
				"group_id": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: "",
				},
				"emoji": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: false,
				},
				"body": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				auth := model.GetAuth(params)
				if auth == nil {
					return nil, errors.New("access denied")
				}

				userId := auth.UserId
				groupId, ok := params.Args["group_id"].(int)

				if !ok {

					return nil, errors.New("invalid group id")

				}
				message := &model.Message{
					UserId:  userId,
					GroupId: int64(groupId),
					Body:    params.Args["body"].(string),
					Emoji:   params.Args["emoji"].(bool),
				}

				fmt.Println("messgae", message)
				result, err := message.Create()

				if err != nil {
					return nil, err
				}

				return result, err

			},
		},
	},
})
