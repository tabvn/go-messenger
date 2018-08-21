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
)

const (
	MysqlConnectUrl = "root:@tcp(127.0.0.1:3306)/messenger?charset=utf8mb4&collation=utf8mb4_unicode_ci"
	IsProduction    = false
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

	_, err := db.InitDatabase(MysqlConnectUrl)
	if err != nil {
		panic(errors.New("can not connect to database"))
	}

}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if !IsProduction {
			// Render GraphIQL
			w.Write(dev.Content)
			return
		}
		content := []byte (`I'm Go!`)
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

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"messenger\": \"v.1.0\"}"))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
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
	mux.HandleFunc("/attachment", model.HandleViewAttachment)
	mux.HandleFunc("/group/avatar", model.HandleViewGroupAvatar)

	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	const PORT = 3007

	fmt.Println("Server is running on port:", PORT)

	err := http.ListenAndServe(":"+strconv.Itoa(PORT), handler)

	if err != nil {
		panic(err)
	}

}
