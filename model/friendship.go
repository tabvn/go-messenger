package model

import (
	"messenger/db"
	"messenger/helper"
	"github.com/graphql-go/graphql"
	"errors"
	"fmt"
	"log"
)

type FriendShip struct {
	Id       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
	Status   bool  `json:"status"`
	Created  int64 `json:"created"`
}

var FriendShipType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"friend_id": &graphql.Field{
				Type: graphql.Int,
			},

			"status": &graphql.Field{
				Type: graphql.Boolean,
			},
			"created": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

func AddFriend(userId, friendId int64) (bool, error) {

	if userId == friendId {
		return false, errors.New("can not add your self")
	}

	status := 0 // default we set 0 as request friend is sent

	count, _ := db.DB.Count("SELECT COUNT(*) FROM friendship WHERE (user_id =? AND friend_id =?) OR (user_id =? AND friend_id =?)", userId, friendId, friendId, userId)

	if count > 0 {
		return ResponseFriendRequest(userId, friendId, true)
	}

	// for now we auto accept from friend.
	q := `INSERT INTO friendship (user_id, friend_id, status, created) VALUES (?, ?, ?, ?), (?, ?, ?, ?)`

	created := helper.GetUnixTimestamp()

	numRows, err := db.DB.InsertMany(q, userId, friendId, status, created, friendId, userId, 1, created)

	if err != nil {

		return false, err
	}
	if numRows == 2 {
		return true, nil
	}

	return false, nil
}

func NotifyFriendAdded(userId, friendId int64) {

	user, e := GetUser(friendId)


	if e == nil && user != nil {
		payload := map[string]interface{}{
			"action": "add_friend",
			"payload": map[string]interface{}{
				"user_id": userId,
				"friend": map[string]interface{}{
					"id":         user.Id,
					"first_name": user.FirstName,
					"last_name":  user.LastName,
					"avatar":     user.Avatar,
					"friend":     true,
					"blocked":    false,
					"status":     UserStatus(user.Online, user.CustomStatus),
				},
			},
		}
		Instance.SendJson(userId, payload)

	}

	// sent to friend
	friend, e := GetUser(userId)

	if e == nil && friend != nil {
		payload := map[string]interface{}{
			"action": "add_friend",
			"payload": map[string]interface{}{
				"user_id": friendId,
				"friend": map[string]interface{}{
					"id":         friend.Id,
					"first_name": friend.FirstName,
					"last_name":  friend.LastName,
					"avatar":     friend.Avatar,
					"friend":     true,
					"blocked":    false,
					"status":     UserStatus(friend.Online, friend.CustomStatus),
				},
			},
		}
		Instance.SendJson(friendId, payload)

	}

}
func ResponseFriendRequest(userId, friendId int64, accepted bool) (bool, error) {

	q := `UPDATE friendship SET status=? WHERE (user_id =? AND friend_id =?) OR (user_id =? AND friend_id =?)`

	status := 1

	if !accepted {
		status = 0
	}

	isUpdate, err := db.DB.Update(q, status, userId, friendId, friendId, userId)
	if err != nil {
		return false, err
	}
	if isUpdate < 1 {
		return false, nil
	}
	if accepted{
		NotifyFriendAdded(userId, friendId)
	}

	return true, nil
}

func UnFriend(userId, friendId int64) (bool, error) {

	q := `DELETE FROM friendship WHERE (user_id =? AND friend_id =?) OR (user_id =? AND friend_id =?)`

	numRows, err := db.DB.DeleteMany(q, userId, friendId, friendId, userId)

	if err != nil {
		log.Println("un friend error", err)
		return false, err
	}
	if numRows == 2 {

		uPayload := map[string]interface{}{
			"action": "un_friend",
			"payload": map[string]interface{}{
				"user_id":   userId,
				"friend_id": friendId,
			},
		}

		fPayload := map[string]interface{}{
			"action": "un_friend",
			"payload": map[string]interface{}{
				"user_id":   friendId,
				"friend_id": userId,
			},
		}

		defer func() {
			Instance.SendJson(userId, uPayload)
			Instance.SendJson(friendId, fPayload)

		}()

		return true, nil
	} else {
		return false, errors.New("an error delete friendship")
	}

	return false, nil
}

func Friends(userId int64, search string, limit, skip int) ([] *User, error) {

	var users []*User

	q := ""

	if search == "" {
		q = `SELECT u.*, f.id, b.id, f.status  
		FROM friendship as f
		INNER JOIN users as u ON f.friend_id = u.id 
		LEFT JOIN blocked as b ON b.author = ? AND b.user = u.id 
		WHERE (SELECT count(*) from friendship where friendship.friend_id = ? AND friendship.status = 1) > 0 AND f.user_id = ? AND f.status = 1 ORDER BY f.created DESC LIMIT ? OFFSET ?`
		rows, err := db.DB.List(q, userId, userId, userId, limit, skip)

		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			user, err := scanUser(rows)

			if err != nil {
				return nil, fmt.Errorf("mysql: could not read row: %v", err)
			}

			user.Password = ""
			users = append(users, user)

		}

	} else {

		// search
		q = `SELECT u.*, f.id, b.id, f.status
			FROM friendship as f
			INNER JOIN users as u ON f.friend_id = u.id 
			LEFT JOIN blocked as b ON b.author = ? AND b.user = u.id 
			WHERE (SELECT count(*) from friendship where friendship.friend_id = ? AND friendship.status = 1) > 0 AND f.user_id = ? AND f.status = 1 AND (u.first_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?) 
			ORDER BY f.created DESC LIMIT ? OFFSET ?`

		search = `%` + search + `%`

		rows, err := db.DB.List(q, userId, userId, userId, search, search, search, limit, skip)

		if err != nil {

			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			user, err := scanUser(rows)

			if err != nil {
				return nil, fmt.Errorf("mysql: could not read row: %v", err)
			}

			user.Password = ""
			users = append(users, user)

		}

	}

	return users, nil

}
