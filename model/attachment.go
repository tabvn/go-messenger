package model

import (
	"github.com/graphql-go/graphql"
	"net/http"
	"messenger/db"
	"database/sql"
	"messenger/config"
)

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

func enableCors(w *http.ResponseWriter) {

	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

func HandleViewAttachment(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	row, err := db.DB.FindOne("SELECT COUNT(*) as count FROM files WHERE name = ?", name)
	if err != nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	var count sql.NullInt64

	if row.Scan(&count) != nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if count.Int64 > 0 {
		http.ServeFile(w, r, config.UploadDir+"/"+name)
		return
	}

	http.Error(w, "access denied", http.StatusForbidden)
	return

}

func HandleViewGroupAvatar(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	row, err := db.DB.FindOne("SELECT COUNT(*) as count FROM groups WHERE avatar = ?", name)
	if err != nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	var count sql.NullInt64

	if row.Scan(&count) != nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if count.Int64 > 0 {
		http.ServeFile(w, r, config.UploadDir+"/"+name)
		return
	}

	http.Error(w, "access denied", http.StatusForbidden)
	return

}

func serveFile(path string, w http.ResponseWriter, r *http.Request) {

}
