package mutation

import (
	"github.com/graphql-go/graphql"
	"messenger/model"
	"errors"
	"messenger/helper"
)

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{

		"createUser": &graphql.Field{
			Type:        model.UserType,
			Description: "Create new user, secret only!",
			Args: graphql.FieldConfigArgument{
				"first_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"uid": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"last_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"avatar": &graphql.ArgumentConfig{
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
					Avatar:    params.Args["avatar"].(string),
					Email:     params.Args["email"].(string),
					Uid:       int64(params.Args["uid"].(int)),
					Password:  params.Args["password"].(string),
				}

				// only allow secret key to create user
				secret := params.Context.Value("secret")

				if secret == nil {

					return nil, errors.New("access denied")
				}

				err := user.Create()

				if err != nil {
					return nil, err
				}

				user.Password = ""

				return user, err

			},
		},
		"createOrUpdateUser": &graphql.Field{
			Type:        model.UserType,
			Description: "Create or update user, Secret only.",
			Args: graphql.FieldConfigArgument{
				"uid": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"first_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"last_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"avatar": &graphql.ArgumentConfig{
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
					Avatar:    params.Args["avatar"].(string),
					Email:     params.Args["email"].(string),
					Uid:       int64(params.Args["uid"].(int)),
					Password:  params.Args["password"].(string),
				}

				// only allow secret key to create user
				secret := params.Context.Value("secret")

				if secret == nil {

					return nil, errors.New("access denied")
				}

				err := user.CreateOrUpdate()

				if err != nil {
					return nil, err
				}

				user.Password = ""

				return user, err

			},
		},
		"requestUserToken": &graphql.Field{
			Type:        model.LoginType,
			Description: "Create or update user, secret only.",
			Args: graphql.FieldConfigArgument{
				"uid": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"first_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"last_name": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"avatar": &graphql.ArgumentConfig{
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
					Avatar:    params.Args["avatar"].(string),
					Email:     params.Args["email"].(string),
					Uid:       int64(params.Args["uid"].(int)),
					Password:  params.Args["password"].(string),
				}

				// only allow secret key to auth user data
				// update or create user. and then create a token
				secret := params.Context.Value("secret")

				if secret == nil {

					return nil, errors.New("access denied")
				}

				token, err := user.RequestUserToken()

				if err != nil {
					return nil, err
				}

				result := map[string]interface{}{
					"id":      token.Id,
					"token":   token.Token,
					"user_id": token.UserId,
					"created": token.Created,
					"user":    user,
				}
				user.Password = ""

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
					Type:         graphql.String,
					DefaultValue: "",
				},
				"avatar": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
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

				secret := params.Context.Value("secret")
				var auth *model.Auth

				if secret == nil {

					auth = model.GetAuth(params)

					if auth == nil {
						return nil, errors.New("access denied")
					} else {
						if int64(id) != auth.UserId {
							// if not super admin only allow self edit
							return nil, errors.New("access denied")
						}
					}

				}

				user := &model.User{
					Id:        int64(id),
					FirstName: params.Args["first_name"].(string),
					LastName:  params.Args["last_name"].(string),
					Avatar:    params.Args["avatar"].(string),
					Email:     params.Args["email"].(string),
					Password:  params.Args["password"].(string),
				}

				err := user.Update()

				if err != nil {
					return nil, err
				}
				return user, err

			},
		},

		"deleteUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete user. Secret only",
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

				secret := params.Context.Value("secret")

				if secret == nil {

					return nil, errors.New("access denied")
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
					"user_id": token.UserId,
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
		"sendMessage": &graphql.Field{
			Type:        model.MessageType,
			Description: "create conversation",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
					Description:  "user_id owner of group",
				},
				"group_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"body": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"emoji": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: false,
				},
				"gif": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"attachments": &graphql.ArgumentConfig{
					Type:         graphql.NewList(graphql.Int),
					DefaultValue: [] int64{},
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				uid := params.Args["user_id"].(int)
				gid := params.Args["group_id"].(int)

				attachments := helper.GetIds(params.Args["attachments"])
				gif := params.Args["gif"].(string)
				body := params.Args["body"].(string)
				emoji := params.Args["emoji"].(bool)

				secret := params.Context.Value("secret")

				userId := int64(uid)
				groupId := int64(gid)

				var auth *model.Auth

				if secret == nil {
					auth = model.GetAuth(params)

					if auth == nil {
						return nil, errors.New("access denied")
					}

					userId = auth.UserId

				}
				message, err := model.CreateMessage(groupId, userId, body, emoji, gif, attachments)

				if err != nil {
					return nil, err
				}

				if message == nil {
					return nil, nil
				}

				return message, nil

			},
		},
		"markAsRead": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Mark message as read",
			Args: graphql.FieldConfigArgument{
				"ids": &graphql.ArgumentConfig{
					Type:         graphql.NewList(graphql.Int),
					DefaultValue: []int{},
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				inputIds, ok := params.Args["ids"]

				if !ok {
					return nil, errors.New("invalid id")
				}

				var auth *model.Auth

				auth = model.GetAuth(params)

				if auth == nil {
					return nil, errors.New("access denied")
				}
				ids := helper.GetIds(inputIds)

				defer func() {
					for _, id := range ids {
						model.MarkMessageAsRead(id, auth.UserId)
					}
				}()

				return true, nil

			},
		},

		"deleteMessage": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete message",
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

				mId := int64(id)

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")
				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {
						// let check if user has perm to delete a message

						canDelete := model.UserCanDeleteMessage(auth.UserId, mId)
						if !canDelete {
							return nil, errors.New("access denied")
						}
					}
				}

				m := model.Message{
					Id: mId,
				}

				result, err := m.Delete()

				if err != nil {
					return nil, errors.New("error")
				}
				return result, err

			},
		},
		"createConversation": &graphql.Field{
			Type:        model.GroupType,
			Description: "create conversation",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
					Description:  "user_id owner of group",
				},
				"participants": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.Int),
				},
				"body": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"emoji": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: false,
				},
				"gif": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"attachments": &graphql.ArgumentConfig{
					Type:         graphql.NewList(graphql.Int),
					DefaultValue: [] int64{},
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				uid := params.Args["user_id"].(int)
				userIdInput := params.Args["participants"]

				userIds := helper.GetIds(userIdInput)

				attachments := helper.GetIds(params.Args["attachments"])
				messageGif := params.Args["gif"].(string)
				messageBody := params.Args["body"].(string)
				messageEmoji := params.Args["emoji"].(bool)

				secret := params.Context.Value("secret")

				userId := int64(uid)

				var auth *model.Auth

				if secret == nil {
					auth = model.GetAuth(params)

					if auth == nil {
						return nil, errors.New("access denied")
					}

					userId = auth.UserId

					if userId > 0 {
						var hasAuthorInList = false

						for _, id := range userIds {
							if id == auth.UserId {
								hasAuthorInList = true
								break
							}
						}

						if !hasAuthorInList {
							userIds = append(userIds, auth.UserId)
						}
					}

				}

				if len(userIds) < 2 {
					return nil, errors.New("must have more than 1 member in conversation")
				}
				group, err := model.CreateConversation(userId, userIds, messageBody, messageGif, messageEmoji, attachments)

				if err != nil {
					return nil, err
				}

				if group == nil {
					return nil, nil
				}
				return group, nil

			},
		},
		"findOrCreateGroup": &graphql.Field{
			Type:        model.GroupType,
			Description: "create group",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
					Description:  "user_id owner of group",
				},
				"participants": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.Int),
				},
				"title": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"avatar": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				uid := params.Args["user_id"].(int)
				userIdInput := params.Args["participants"]
				title := params.Args["title"].(string)
				avatar := params.Args["avatar"].(string)

				userIds := helper.GetIds(userIdInput)

				secret := params.Context.Value("secret")

				userId := int64(uid)

				var auth *model.Auth

				if secret == nil {
					auth = model.GetAuth(params)

					if auth == nil {
						return nil, errors.New("access denied")
					}

					userId = auth.UserId

					if userId > 0 {
						var hasAuthorInList = false

						for _, id := range userIds {
							if id == auth.UserId {
								hasAuthorInList = true
								break
							}
						}

						if !hasAuthorInList {
							userIds = append(userIds, auth.UserId)
						}
					}

				}
				gid, err := model.FindOrCreateGroup(userId, userIds, title, avatar)

				if err != nil {
					return nil, errors.New("error")
				}

				if gid == 0 {
					return nil, nil
				}

				group, loadErr := model.LoadGroup(gid, userId)

				if loadErr != nil {
					return nil, loadErr
				}

				return group, nil

			},
		},
		"joinGroup": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "add user to group chat",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"group_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user_id"].(int)
				groupId, groupOk := params.Args["group_id"].(int)

				if !ok {
					return nil, errors.New("invalid id")
				}
				if ! groupOk {

					return nil, errors.New("invalid group id")
				}

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")

				uid := int64(userId)
				gid := int64(groupId)

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {
						// let check if user has perm to delete a message
						canJoin := model.CanJoinGroup(auth.UserId, uid, gid)

						if !canJoin {
							return nil, errors.New("access denied")
						}
					}
				}

				result := model.JoinGroup(uid, gid)

				return result, nil

			},
		},
		"leftGroup": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "left group chat",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"group_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user_id"].(int)
				groupId, groupOk := params.Args["group_id"].(int)

				if !ok {
					return nil, errors.New("invalid id")
				}
				if ! groupOk {

					return nil, errors.New("invalid group id")
				}

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")

				uid := int64(userId)
				gid := int64(groupId)

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {

						// if not super admin only accept userId from auth request

						uid = auth.UserId

					}
				}

				_, err := model.LeftGroup(uid, gid)

				if err != nil {
					return false, errors.New("not found")
				}

				return true, nil

			},
		},
		"addFriend": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "add friend",
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Description:  "user_id, only allow set user_id for secret, other wise take user_id from auth",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"friend": &graphql.ArgumentConfig{
					Description:  "Friend user_id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user"].(int)

				friendId, fok := params.Args["friend"].(int)

				if !fok {
					return nil, errors.New("invalid friend user_id")
				}
				if !ok {
					return nil, errors.New("invalid user_id")
				}

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")

				uid := int64(userId)
				friend := int64(friendId)

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {

						// only take user_id from auth

						uid = auth.UserId

					}
				}

				result, err := model.AddFriend(uid, friend)

				if err != nil {
					return false, err
				}

				return result, err

			},
		},
		"blockUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "remove friendship",
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Description:  "block user",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"friend": &graphql.ArgumentConfig{
					Description:  "Friend user_id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user"].(int)

				friendId, fok := params.Args["friend"].(int)

				if !fok {
					return nil, errors.New("invalid friend user_id")
				}
				if !ok {
					return nil, errors.New("invalid user_id")
				}

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")

				uid := int64(userId)
				friend := int64(friendId)

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {

						// only take user_id from auth
						uid = auth.UserId

					}
				}
				result, err := model.BlockUser(uid, friend)

				if err != nil {
					return false, err
				}

				return result, err

			},
		},
		"unBlockUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "remove friendship",
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Description:  "block user",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"friend": &graphql.ArgumentConfig{
					Description:  "Friend user_id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user"].(int)

				friendId, fok := params.Args["friend"].(int)

				if !fok {
					return nil, errors.New("invalid friend user_id")
				}
				if !ok {
					return nil, errors.New("invalid user_id")
				}

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")

				uid := int64(userId)
				friend := int64(friendId)

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {

						// only take user_id from auth
						uid = auth.UserId

					}
				}

				result, err := model.UnBlockUser(uid, friend)

				if err != nil {
					return false, err
				}

				return result, err

			},
		},
		"unFriend": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "remove friendship",
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Description:  "user_id, only allow set user_id for secret, other wise take user_id from auth",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"friend": &graphql.ArgumentConfig{
					Description:  "Friend user_id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user"].(int)

				friendId, fok := params.Args["friend"].(int)

				if !fok {
					return nil, errors.New("invalid friend user_id")
				}
				if !ok {
					return nil, errors.New("invalid user_id")
				}

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")

				uid := int64(userId)
				friend := int64(friendId)

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {

						// only take user_id from auth

						uid = auth.UserId

					}
				}
				result, err := model.UnFriend(uid, friend)

				if err != nil {
					return false, err
				}

				return result, err

			},
		},
	},
})
