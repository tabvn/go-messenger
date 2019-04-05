package main

import (
	"net/http"
	"encoding/json"
	"errors"
	"context"
	"messenger/dev"
	"messenger/schema"
	"messenger/model"
	"messenger/db"
	"fmt"
	"strconv"
	"messenger/upload"
	"github.com/rs/cors"
	"messenger/config"
	"database/sql"
)

type params struct {
	Query         string      `json:"query"`
	OperationName string      `json:"operationName,omitempty"`
	Variables     interface{} `json:"variables,omitempty"`
}

func getBodyFromRequest(r *http.Request) (*params, error) {
	p := &params{
		Variables: nil,
	}

	if r.Method == "POST" {

		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func Setup() {

	_, err := db.InitDatabase(config.MysqlConnectUrl)
	if err != nil {
		panic(errors.New("can not connect to database"))
	}

	db.DB.Update("UPDATE users SET online =?", false)

}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if !config.Production {
			// Render GraphIQL
			w.Write(dev.Content)
			return
		}
		content := []byte (`v.1.0`)
		w.Write(content)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	p, err := getBodyFromRequest(r)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)

		return
	}

	var auth = r.Header.Get("Authorization")

	if len(auth) == 0 {
		auth = r.URL.Query().Get("auth")
	}

	cookie := model.GetCookie(r, "token")
	if cookie != auth {
		model.SetCookie(w, "token", auth)
	}

	isSecret := model.CheckSecret(auth)

	var ctx context.Context

	if isSecret {

		ctx = context.WithValue(context.Background(), "secret", isSecret)
	} else {
		authentication, _ := model.VerifyToken(auth)
		ctx = context.WithValue(context.Background(), "auth", authentication)
	}

	result := schema.ExecuteQuery(ctx, p.Query, p.OperationName, schema.Schema)

	json.NewEncoder(w).Encode(result)
}

func enableCors(w *http.ResponseWriter) {

	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

func handleViewAttachment(w http.ResponseWriter, r *http.Request) {

	cookie := model.GetCookie(r, "token")

	if cookie == "" {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}
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
func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"version\": \"v.1.0\"}"))
	})

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("User-agent: *\nDisallow: /"))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*.pantheonsite.io","http://localhost:3000","http://drupal7.test","http://127.0.0.1:3000", "https://*.addictionrecovery.com", "http://*.addictionrecovery.com"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Origin",
			"Origin",
			"Access-Control-Request-Headers",
			"Access-Control-Request-Method",
			"Connection",
		},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	// Use default options
	handler := c.Handler(mux)

	Setup()
	// Router api graphQL handler
	mux.HandleFunc("/api", graphqlHandler)
	mux.HandleFunc("/ws", model.WebSocketHandler)
	mux.HandleFunc("/upload", upload.HandleFileUpload)
	mux.HandleFunc("/uploads", upload.HandleMultiUpload)
	mux.HandleFunc("/attachment", handleViewAttachment)
	mux.HandleFunc("/group/avatar", model.HandleViewGroupAvatar)

	fs := http.FileServer(http.Dir(config.PublicDir))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	fmt.Println("Server is running on port:", config.Port)

	panic(http.ListenAndServe(":"+strconv.Itoa(config.Port), handler))

}
