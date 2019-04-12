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
				"published": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 1,
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
					Published: int64(params.Args["published"].(int)),
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
				"published": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 1,
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
					Published: int64(params.Args["published"].(int)),
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
				"published": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 1,
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
					Published: int64(params.Args["published"].(int)),
				}

				err := user.Update()

				if err != nil {
					return nil, err
				}
				return user, err

			},
		},

		"updateUserStatus": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "update user status",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"status": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.String),
					DefaultValue: "online",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user_id"].(int)
				status, k := params.Args["status"].(string)

				if !ok {
					return nil, errors.New("invalid user_id")
				}
				if !k {

					return nil, errors.New("invalid status")
				}

				var auth *model.Auth
				uid := int64(userId)

				secret := params.Context.Value("secret")

				if secret == nil {

					auth = model.GetAuth(params)
					uid = auth.UserId

					if auth == nil {
						return nil, errors.New("access denied")
					}

				}

				if uid == 0 {
					return nil, errors.New("invalid user_id")
				}

				isOK := model.UpdateUserStatus(uid, true, status)

				return isOK, nil

			},
		},
		"deleteUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete user. Secret only",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"uid": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				id, ok := params.Args["id"].(int)
				uid, k := params.Args["uid"].(int)

				if !ok {
					return nil, errors.New("invalid id")
				}

				if !k {
					return nil, errors.New("invalid uid");
				}

				secret := params.Context.Value("secret")

				if secret == nil {

					return nil, errors.New("access denied")
				}

				if id == 0 && uid == 0 {
					return nil, errors.New("invalid user_id or uid")
				}
				result := false
				var err error

				if id > 0 {
					user := model.User{
						Id: int64(id),
					}

					result, err = user.Delete()
				} else {

					result, err = model.DeleteUserBy(int64(uid))
				}

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
			Type:        graphql.Boolean,
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

				return success, err

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

		"responseInvite": &graphql.Field{
			Type:        graphql.Boolean,
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
				"accept": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: true,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				uid := params.Args["user_id"].(int)
				gid := params.Args["group_id"].(int)
				accept := params.Args["accept"].(bool)
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

				err := model.ResponseInvite(userId, groupId, accept)

				return true, err

			},
		},

		"archiveGroup": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Mark message as read",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"group_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, k := params.Args["user_id"].(int)
				if !k {
					return nil, errors.New("invalid user_id")
				}

				groupId, ok := params.Args["group_id"].(int)

				if !ok {
					return nil, errors.New("invalid group_id")
				}

				gid := int64(groupId)
				uid := int64(userId)

				var auth *model.Auth

				secret := params.Context.Value("secret")

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					}

					uid = auth.UserId

				}

				if uid == 0 || gid == 0 {
					return nil, errors.New("invalid user_id or group_id")
				}

				b := model.ArchiveGroup(uid, gid)

				return b, nil

			},
		},

		"markAsRead": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Mark message as read",
			Args: graphql.FieldConfigArgument{
				"group_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"ids": &graphql.ArgumentConfig{
					Type:         graphql.NewList(graphql.Int),
					DefaultValue: []int{},
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				var auth *model.Auth

				auth = model.GetAuth(params)

				if auth == nil {
					return nil, errors.New("access denied")
				}

				gid, ok := params.Args["group_id"].(int)

				if !ok {
					return nil, errors.New("invalid group_id")
				}

				groupId := int64(gid)
				if groupId > 0 {
					err := model.MarkAsReadByGroup(groupId, auth.UserId)
					if err != nil {
						return false, errors.New("an error")
					}
					return true, nil
				}

				inputIds, ok := params.Args["ids"]

				if !ok {
					return nil, errors.New("invalid id")
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

		"updateMessage": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "update message",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"body": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.String),
					DefaultValue: "",
				},
				"emoji": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: false,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				body := params.Args["body"].(string)
				emoji := params.Args["emoji"].(bool)
				uid, k := params.Args["user_id"].(int)
				if !k {
					return nil, errors.New("invalid user_id")
				}
				id, ok := params.Args["id"].(int)
				if !ok {
					return nil, errors.New("invalid id")
				}

				messageId := int64(id)
				userId := int64(uid)

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")
				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					}

					userId = auth.UserId
				}

				bool := model.UpdateMessage(userId, messageId, body, emoji)

				return bool, nil

			},
		},

		"deleteMessage": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete message",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				uid, k := params.Args["user_id"].(int)
				if !k {
					return nil, errors.New("invalid user_id")
				}
				id, ok := params.Args["id"].(int)
				if !ok {
					return nil, errors.New("invalid id")
				}

				messageId := int64(id)
				userId := int64(uid)

				var auth *model.Auth

				// allow super or authenticated user
				secret := params.Context.Value("secret")
				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					}

					userId = auth.UserId
				}

				bool := model.DeleteMessage(userId, messageId)

				return bool, nil

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
				"title": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"avatar": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
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

				groupTitle := params.Args["title"].(string)
				groupAvatar := params.Args["avatar"].(string)

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
				group, err := model.CreateConversation(userId, userIds, messageBody, messageGif, messageEmoji, attachments, groupTitle, groupAvatar)

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
		"updateGroup": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "update group",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
					Description:  "group id",
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

				id := params.Args["id"].(int)
				title := params.Args["title"].(string)
				avatar := params.Args["avatar"].(string)

				secret := params.Context.Value("secret")

				groupId := int64(id)
				var auth *model.Auth
				var userId int64

				if secret == nil {
					auth = model.GetAuth(params)

					if auth == nil {
						return false, errors.New("access denied")
					}

					userId = auth.UserId

					if !model.IsMemberOfGroup(userId, groupId) {
						return nil, errors.New("access denied")
					}

				}

				return model.UpdateGroup(groupId, title, avatar), nil

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

				var addByUserId int64

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {
						// let check if user has perm to delete a message
						canJoin := model.CanJoinGroup(auth.UserId, uid, gid)

						addByUserId = auth.UserId

						if !canJoin {
							return nil, errors.New("access denied")
						}
					}
				}

				result := model.JoinGroup(uid, gid, addByUserId)

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
						if uid == 0 {
							uid = auth.UserId
						}

						// check if user has perm to remove

						if auth.UserId != uid {
							// let check if user is has perm to remove a user from group

							canRemove := model.CanDeleteMember(auth.UserId, uid, gid)
							if !canRemove {
								return nil, errors.New("can not remove user")
							}
						}
					}
				}

				_, err := model.LeftGroup(uid, gid)

				if err != nil {
					return false, errors.New("not found")
				}

				return true, nil

			},
		},
		"removeGroupUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "remove user from group chat",
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

				var removeByUserId int64

				if secret == nil {
					auth = model.GetAuth(params)
					if auth == nil {
						return nil, errors.New("access denied")
					} else {

						// check if user has perm to remove

						removeByUserId = auth.UserId
						canRemove := model.CanDeleteMember(auth.UserId, uid, gid)
						if !canRemove {
							return nil, errors.New("can not remove user")
						}
					}
				}

				err := model.RemoveGroupUser(removeByUserId, uid, gid)

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
					Type:         graphql.Int,
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
		"ResponseFriendRequest": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "remove friendship",
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Description:  "block user",
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"friend": &graphql.ArgumentConfig{
					Description:  "Friend user_id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"accept": &graphql.ArgumentConfig{
					Description:  "Accept or deny",
					Type:         graphql.Boolean,
					DefaultValue: true,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, ok := params.Args["user"].(int)

				friendId, fok := params.Args["friend"].(int)
				accepted := params.Args["accept"].(bool)

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

				result, err := model.ResponseFriendRequest(uid, friend, accepted)

				if err != nil {
					return false, errors.New("an error response friend request")
				}

				return result, err

			},
		},
		"userTyping": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "User typing",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Description:  "user id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"group_id": &graphql.ArgumentConfig{
					Description:  "group_id",
					Type:         graphql.NewNonNull(graphql.Int),
					DefaultValue: 0,
				},
				"is_typing": &graphql.ArgumentConfig{
					Description:  "is typing",
					Type:         graphql.Boolean,
					DefaultValue: true,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userId, _ := params.Args["user_id"].(int)

				groupId, _ := params.Args["group_id"].(int)
				isTyping := params.Args["is_typing"].(bool)


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

						// only take user_id from auth
						uid = auth.UserId

					}
				}

				result := model.UserIsTyping(gid, uid, isTyping)
				return result, nil

			},
		},
		"blockUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "remove friendship",
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Description:  "block user",
					Type:         graphql.Int,
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
					return false, errors.New("an error block user")
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
					Type:         graphql.Int,
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
					return false, errors.New("an error")
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
