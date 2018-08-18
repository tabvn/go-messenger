package upload

import (
	"crypto/rand"
	"net/http"
	"fmt"
	"io/ioutil"
	"mime"
	"path/filepath"
	"os"
	"messenger/model"
	"encoding/json"
)

const maxUploadSize = 1024 * 1024

func renderError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}
func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func upload(userId int64, w http.ResponseWriter, r *http.Request) () {

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}

	// parse and validate file and post parameters
	file, handler, err := r.FormFile("file")
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	filetype := http.DetectContentType(fileBytes)

	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return
	}
	fileName := randToken(12)
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return
	}
	name := fileName + fileEndings[0]
	newPath := filepath.Join("./storage", name)
	fmt.Printf("FileType: %s, File: %s\n", filetype, newPath)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close() // idempotent, okay to call twice

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}

	fileSaved, err := model.SaveFile(userId, name, handler.Filename, filetype, handler.Size)
	if err != nil {

		http.Error(w, "unable save file", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileSaved)

	return

}
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {

	var auth = r.Header.Get("Authorization")

	if len(auth) == 0 {
		auth = r.URL.Query().Get("auth")
	}

	var userId int64

	authentication, err := model.VerifyToken(auth)

	if err != nil || authentication == nil {
		http.Error(w, "access denied", http.StatusForbidden)

		return
	}

	if authentication != nil {
		userId = authentication.UserId
	}

	if r.Method == "GET" {
		http.Error(w, "error", http.StatusForbidden)
		return

	} else {

		upload(userId, w, r)

	}
}
