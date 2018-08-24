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
	"io"
	"messenger/sanitize"
	"messenger/config"
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

func HandleMultiUpload(w http.ResponseWriter, r *http.Request) {

	var auth = r.Header.Get("Authorization")

	if len(auth) == 0 {
		auth = r.URL.Query().Get("auth")
	}

	authentication, err := model.VerifyToken(auth)

	if err != nil || authentication == nil {
		http.Error(w, "access denied", http.StatusForbidden)

		return
	}

	e := r.ParseMultipartForm(1048576)
	if e != nil {

		renderError(w, "ERROR", http.StatusBadRequest)

		return
	}

	files := r.MultipartForm.File["files"] //request.MultipartForm has multipart.Form and multipart.Form.File[] has type FileHeader, not File. So, you have to Open() this.

	var uploaded [] model.File

	for _, file := range files {
		f, err := file.Open()
		defer f.Close()
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		fileType := http.DetectContentType(fileBytes)

		fileName := randToken(12)

		name := fileName + "_" + sanitize.Name(file.Filename)
		newPath := filepath.Join(config.UploadDir, name)

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

		fileSaved, err := model.SaveFile(authentication.UserId, name, file.Filename, fileType, file.Size)
		if err != nil {
			http.Error(w, "unable save file", http.StatusForbidden)
			return
		}
		uploaded = append(uploaded, *fileSaved)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploaded)

}

func upload(userId int64, w http.ResponseWriter, r *http.Request) () {

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
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

	fileType := http.DetectContentType(fileBytes)

	fileName := randToken(12)
	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return
	}
	name := fileName + fileEndings[0]
	newPath := filepath.Join(config.UploadDir, name)
	fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

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

	fileSaved, err := model.SaveFile(userId, name, handler.Filename, fileType, handler.Size)
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
