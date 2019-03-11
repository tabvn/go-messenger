package model

import (
	"time"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"messenger/helper"
	"errors"
	"messenger/db"
	"github.com/satori/go.uuid"
	"fmt"
	"database/sql"
	"strconv"
	"messenger/config"
)

const (
	ONLINE  = "online"
	BUSY    = "busy"
	AWAY    = "away"
	OFFLINE = "offline"
)

type User struct {
	Id                int64  `json:"id"`
	Uid               int64  `json:"uid"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	Avatar            string `json:"avatar"`
	Online            bool   `json:"online"`
	CustomStatus      string `json:"custom_status"`
	Status            string `json:"status"`
	Location          string `json:"location"`
	Work              string `json:"work"`
	School            string `json:"school"`
	About             string `json:"about"`
	Created           int64  `json:"created"`
	Updated           int64  `json:"updated"`
	Friend            bool   `json:"friend"`
	FriendRequestSent bool   `json:"friend_request_sent"`
	Blocked           bool   `json:"blocked"`
	Published         int64  `json:"published"`
}

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"uid": &graphql.Field{
				Type: graphql.Int,
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
			"avatar": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"location": &graphql.Field{
				Type: graphql.String,
			},
			"work": &graphql.Field{
				Type: graphql.String,
			},
			"school": &graphql.Field{
				Type: graphql.String,
			},
			"about": &graphql.Field{
				Type: graphql.String,
			},
			"created": &graphql.Field{
				Type: graphql.Int,
			},
			"updated": &graphql.Field{
				Type: graphql.Int,
			},
			"friend": &graphql.Field{
				Type: graphql.Boolean,
			},
			"blocked": &graphql.Field{
				Type: graphql.Boolean,
			},
			"friend_request_sent": &graphql.Field{
				Type: graphql.Boolean,
			},
			"published": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var LoginType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "login",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
			"user_id": &graphql.Field{
				Type: graphql.ID,
			},
			"created": &graphql.Field{
				Type: graphql.Int,
			},
			"user": &graphql.Field{
				Type: UserType,
			},
		},
	},
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) Create() (error) {

	validateError := u.validateCreate()

	if validateError != nil {

		return validateError
	}

	// generate password
	password, e := HashPassword(u.Password)
	u.Password = password

	if e != nil {
		return e
	}

	query := `INSERT INTO users (uid, first_name, last_name, email, avatar , password, published,created, updated) VALUES (?,?, ?, ?, ?, ?, ?, ?, ?)`
	currentTime := time.Now()
	u.Created = currentTime.Unix()
	u.Updated = currentTime.Unix()

	result, err := db.DB.Insert(query, u.Uid, u.FirstName, u.LastName, u.Email, u.Avatar, u.Password, u.Published, u.Created, u.Updated)

	if err != nil {

		return err
	}

	u.Id = result

	u.Password = ""

	return err
}

func (u *User) Update() (error) {

	currentTime := time.Now()
	u.Updated = currentTime.Unix()

	if u.Password == "" {
		query := `UPDATE users SET first_name=?, last_name=?, email=?, avatar =?, updated=?, published=? WHERE id = ?`
		_, err := db.DB.Update(query, u.FirstName, u.LastName, u.Email, u.Avatar, u.Updated, u.Published, u.Id)

		if err != nil {
			return err
		}
	} else {
		query := `UPDATE users SET first_name=?, last_name=?, email=?, password=?, avatar =?, updated=?, published=? WHERE id = ?`
		password, err := HashPassword(u.Password)
		if err != nil {
			return err
		}
		_, updateErr := db.DB.Update(query, u.FirstName, u.LastName, u.Email, password, u.Avatar, u.Updated, u.Published, u.Id)

		if updateErr != nil {
			return err
		}
	}

	u.Password = ""

	return nil
}

func (u *User) CreateOrUpdate() (error) {

	// let find if user exist
	findQuery := `SELECT id, COUNT(*) as c FROM users WHERE uid = ?`

	row, err := db.DB.FindOne(findQuery, u.Uid)

	if err != nil {
		return err
	}
	var (
		id    sql.NullInt64
		count int64
	)

	scanErr := row.Scan(&id, &count)

	if scanErr != nil && scanErr != sql.ErrNoRows {
		return scanErr
	}

	if count > 0 {
		// found user exist do update
		u.Id = id.Int64

		updateErr := u.Update()
		if updateErr != nil {
			return updateErr
		}

		return nil

	} else {

		// do create
		if u.Password == "" {
			u.Password = helper.GenerateID()
		}
		createErr := u.Create()

		if createErr != nil {
			return createErr
		}
	}

	return nil
}

func (u *User) RequestUserToken() (*Token, error) {

	err := u.CreateOrUpdate()

	if err != nil {
		return nil, err
	}

	t := &Token{
		Id:      0,
		UserId:  u.Id,
		Token:   uuid.Must(uuid.NewV4()).String(),
		Created: helper.GetUnixTimestamp(),
	}

	r, tokenErr := t.Create()

	if tokenErr != nil {
		return nil, tokenErr
	}

	return r, nil
}

func UserStatus(online bool, customStatus string) (string) {

	var onlineStatus string

	if online {
		onlineStatus = ONLINE
	} else {
		onlineStatus = OFFLINE
	}

	if online && (customStatus == ONLINE || customStatus == OFFLINE || customStatus == BUSY || customStatus == AWAY) {
		onlineStatus = customStatus
	}

	return onlineStatus
}
func scanUser(s db.RowScanner) (*User, error) {

	var (
		id           int64
		uid          int64
		firstName    sql.NullString
		lastName     sql.NullString
		email        sql.NullString
		password     sql.NullString
		avatar       sql.NullString
		online       sql.NullBool
		customStatus sql.NullString
		location     sql.NullString
		work         sql.NullString
		school       sql.NullString
		about        sql.NullString
		created      sql.NullInt64
		updated      sql.NullInt64
		published    sql.NullInt64
		friendship   sql.NullInt64
		blocked      sql.NullInt64
		friendStatus sql.NullInt64
	)

	if err := s.Scan(&id, &uid, &firstName, &lastName, &email, &password, &avatar,
		&online, &customStatus, &location, &work, &school, &about, &created, &updated, &published, &friendship, &blocked, &friendStatus);
		err != nil {

		return nil, err
	}

	onlineStatus := UserStatus(online.Bool, customStatus.String)

	var isFriend = false
	var isBlocked = false
	var isRequestFriendSent = false

	if friendship.Int64 > 0 {
		isFriend = true
	}
	if blocked.Int64 > 0 {
		isBlocked = true
	}
	if friendStatus.Valid && friendStatus.Int64 == 0 {
		isRequestFriendSent = true
	}
	user := &User{
		Id:                id,
		Uid:               uid,
		FirstName:         firstName.String,
		LastName:          lastName.String,
		Email:             email.String,
		Password:          password.String,
		Avatar:            avatar.String,
		Online:            online.Bool,
		Status:            onlineStatus,
		CustomStatus:      customStatus.String,
		Location:          location.String,
		Work:              work.String,
		School:            school.String,
		About:             about.String,
		Created:           created.Int64,
		Updated:           updated.Int64,
		Friend:            isFriend,
		Blocked:           isBlocked,
		FriendRequestSent: isRequestFriendSent,
		Published:         published.Int64,
	}

	if published.Int64 == 0 || (published.Int64 == 2 && !isFriend) {
		user.FirstName = "Anonymous"
		user.LastName = ""
		user.Avatar = config.PrivateAvatar
	}

	return user, nil
}

func (u *User) Load() (*User, error) {

	//count(id), count(created) is fake scan for blocked and friend

	var row *sql.Row
	var err error

	if u.Uid > 0 {
		row, err = db.DB.FindOne(`SELECT u.*, count(id), count(created), count(updated)  FROM users AS u  WHERE uid = ?`, u.Uid)

	} else {
		row, err = db.DB.FindOne(`SELECT u.*, count(id), count(created) , count(updated) FROM users AS u  WHERE id = ?`, u.Id)
	}

	if err != nil {
		return nil, err
	}

	user, err := scanUser(row)

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, err
}
func GetUser(id int64) (*User, error) {

	row, err := db.DB.FindOne(`SELECT u.*, count(id), count(created),count(updated)  FROM users AS u  WHERE id = ?`, id)

	if err != nil {
		return nil, err
	}

	user, err := scanUser(row)

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, err

}

func (u *User) Delete() (bool, error) {

	_, err := db.DB.Delete("DELETE FROM users where id=?", u.Id)

	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteUserBy(uid int64) (bool, error) {

	_, err := db.DB.Delete("DELETE FROM users where uid=?", uid)

	if err != nil {
		return false, err
	}

	return true, nil
}

func VerifyToken(token string) (*Auth, error) {

	if token == "" {
		return nil, nil
	}

	row, err := db.DB.FindOne("SELECT * FROM tokens WHERE token=?", token)

	if err != nil {
		writeToLog("Token error "+err.Error(), "verify_token")
		return nil, errors.New("invalid token")
	}

	t, err := scanToken(row)

	if err != nil {

		writeToLog("scan token error "+err.Error(), "verify_token")

		return nil, errors.New("invalid token")
	}

	var user = &User{Id: t.UserId}

	u, err := user.Load()

	if err != nil {
		writeToLog("could not load user with token "+err.Error(), "verify_token")

		return nil, errors.New("invalid token")
	}

	auth := &Auth{t, u}

	return auth, err

}
func (u *User) validateCreate() (error) {

	var err error = nil

	// Email validation
	if u.Email == "" {
		err = errors.New("email is required")
		return err
	}

	u.Email = strings.ToLower(u.Email)
	err = helper.ValidateEmail(u.Email)

	if err != nil {
		return err
	}

	count, countErr := db.DB.Count("SELECT COUNT(*) FROM users WHERE email=?", u.Email)

	if countErr != nil {
		return errors.New("unable validate email")
	}
	if count > 0 {
		return errors.New("email already exist")
	}

	// trim space
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	// Password validation
	if u.Password == "" {
		err = errors.New("password is required")
		return err
	}

	if len(u.Password) < 6 {
		err = errors.New("password must be of minimum 6 characters length")
		return err
	}

	return err
}

func LoginUser(email string, password string) (*Token, *User, error) {

	row := db.DB.QueryRow("SELECT * FROM users WHERE email=?", email)

	user, err := scanUser(row)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("login failure")
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, nil, errors.New("login failure")
	}

	currentTime := time.Now()

	t := &Token{
		Id:      0,
		UserId:  user.Id,
		Token:   uuid.Must(uuid.NewV4()).String(),
		Created: currentTime.Unix(),
	}

	r, createTokenErr := t.Create()

	if createTokenErr != nil {
		return nil, nil, createTokenErr
	}

	return r, user, nil

}

func LogoutUser(token string) (bool, error) {

	var success = false

	_, err := db.DB.Delete("DELETE FROM tokens where token =?", token)

	if err != nil {

		return false, err
	} else {
		success = true
	}

	return success, nil
}

func GetBlockedUsers(userId int64, limit int, skip int) ([]*User, error) {

	q := `SELECT u.*, f.id, b.id, f.status FROM users as u
			LEFT JOIN friendship as f ON  f.friend_id = u.id AND f.user_id =? AND f.status = 1 
			INNER JOIN blocked as b ON b.author =? AND b.user = u.id
			ORDER BY created DESC LIMIT ? OFFSET ?`

	rows, err := db.DB.List(q, userId, userId, limit, skip)

	if err != nil {
		return nil, err
	}

	var users []*User

	defer rows.Close()

	for rows.Next() {
		user, err := scanUser(rows)

		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}

		user.Password = ""
		users = append(users, user)

	}

	return users, nil

}
func Users(userId int64, search string, limit int, skip int) ([]*User, error) {

	var rows *sql.Rows
	var err error

	if search == "" {

		q := `SELECT u.*, f.id, b.id, f.status FROM users as u 
			LEFT JOIN friendship as f ON  f.friend_id = u.id AND f.user_id =? AND f.status != 2 
			LEFT JOIN blocked as b ON b.author =? AND b.user = u.id WHERE 
			u.id NOT IN (SELECT user FROM blocked WHERE author =? AND user = u.id)
			AND u.id NOT IN (SELECT author FROM blocked WHERE author = u.id AND user = ?) 
			ORDER BY created DESC LIMIT ? OFFSET ?`

		rows, err = db.DB.List(q, userId, userId, userId, userId, limit, skip)

	} else {

		like := "%" + search + "%"
		q := `SELECT u.*, f.id, b.id, f.status FROM users as u 
			LEFT JOIN friendship as f ON  f.friend_id = u.id AND f.user_id =? AND f.status != 2 
			LEFT JOIN blocked as b ON b.author =? AND b.user = u.id 
			WHERE 
			u.id NOT IN (SELECT user FROM blocked WHERE author =? AND user = u.id)
			AND u.id NOT IN (SELECT author FROM blocked WHERE author = u.id AND user = ?)
			AND (u.first_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?)
			ORDER BY created DESC LIMIT ? OFFSET ?`
		rows, err = db.DB.List(q, userId, userId, userId, userId, like, like, like, limit, skip)

	}

	if err != nil {
		return nil, err
	}

	var users []*User

	defer rows.Close()

	for rows.Next() {
		user, err := scanUser(rows)

		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}

		user.Password = ""
		users = append(users, user)

	}

	return users, nil
}

func CountUsers() (int, error) {

	count, err := db.DB.Count("SELECT COUNT(*) FROM users")

	if err != nil {
		return 0, err
	}

	return count, nil
}

func BlockUser(userId, friendId int64) (bool, error) {

	if userId == friendId {
		return false, errors.New("can not blocked your self")
	}

	q := `INSERT INTO blocked (author, user) VALUES (?, ?)`

	id, err := db.DB.Insert(q, userId, friendId)

	if err != nil {
		return false, err

	}
	if id > 0 {
		return true, nil
	}

	return false, nil
}

func UnBlockUser(userId, friendId int64) (bool, error) {

	if userId == friendId {
		return false, errors.New("can not un blocked your self")
	}

	q := `DELETE FROM blocked WHERE author =? AND user =?`

	_, err := db.DB.Delete(q, userId, friendId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func FindUserToNotify(userId int64) ([] int64) {

	var list []int64

	q := `SELECT a.id FROM 
	(SELECT DISTINCT(u.id) FROM users as u INNER JOIN members as m ON u.id = m.user_id AND u.online = 1 AND u.id !=? WHERE m.group_id IN (SELECT groups.id FROM members INNER JOIN groups ON members.group_id = groups.id AND members.user_id =?) ) as a
	LEFT JOIN (SELECT DISTINCT(u.id) from friendship as f INNER JOIN users as u ON u.id = f.friend_id AND f.user_id=? WHERE u.online = true) as b ON b.id = a.id WHERE a.id NOT IN (SELECT b.user FROM blocked AS b WHERE b.author =?)`
	rows, err := db.DB.List(q, userId, userId, userId, userId)

	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var id sql.NullInt64

		if rows.Scan(&id) != nil {
			return nil
		}

		if id.Int64 > 0 {
			list = append(list, id.Int64)
		}

	}

	return list
}

func Notify(userIds [] int64, message []byte) {

	for _, id := range userIds {

		Instance.Send(id, message)
	}
}
func UpdateUserStatus(userId int64, online bool, status string) (bool) {

	var err error

	if status != "" {
		query := `UPDATE users SET online=?, custom_status=? WHERE id = ?`
		_, err = db.DB.Update(query, online, status, userId)

	} else {
		query := `UPDATE users SET online=? WHERE id = ?`
		_, err = db.DB.Update(query, online, userId)
	}

	if err == nil {

		statusQuery := "SELECT custom_status from users WHERE id =?"

		row, e := db.DB.FindOne(statusQuery, userId)

		if e == nil {

			var statusScan sql.NullString
			if row.Scan(&statusScan) == nil {
				status = statusScan.String
			}

		}

		realStatus := UserStatus(online, status)

		message := []byte(`{"action": "user_status", "payload": {"user_id": ` + strconv.Itoa(int(userId)) + `, "status": "` + realStatus + `"}}`)

		var userIds []int64

		userIds = FindUserToNotify(userId)
		userIds = append(userIds, userId)

		defer Notify(userIds, message)

		return true

	}

	return false

}
