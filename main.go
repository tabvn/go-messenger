package main

import (
	"net/http"
	"encoding/json"
	"errors"
	"context"
	"messenger/dev"
	"messenger/schema"
	"messenger/model"
	"messenger/pubsub"
	"messenger/db"
	"fmt"
	"strconv"
	"messenger/upload"
)

const (
	MysqlConnectUrl = "root:@/messenger"
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

func enableCors(w *http.ResponseWriter) {


	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	if (*r).Method == "OPTIONS" {
		return
	}
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

	Setup()
	// Router api graphQL handler
	http.HandleFunc("/api", graphqlHandler)
	http.HandleFunc("/ws", pubsub.WebSocketHandler)
	http.HandleFunc("/files", upload.HandleFileUpload)

	const PORT = 3007

	fmt.Println("Server is running on port:", PORT)

	err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	if err != nil {
		panic(err)
	}

}
