package upload

import (
	"net/http"
	"os"
	"io"
	"encoding/json"
	"messenger/model"
	"messenger/helper"
)

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

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")

		if err != nil || file == nil {
			http.Error(w, "upload error", http.StatusBadRequest)

			return
		}

		defer file.Close()
		name := helper.GenerateID() + "_" + handler.Filename

		f, err := os.OpenFile("./storage/"+name, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, "unable save file", http.StatusForbidden)
			return
		}

		if file != nil {

			fileSaved, err := model.SaveFile(userId, name, handler.Filename, handler.Header.Get("Content-Type"), handler.Size)
			if err != nil {

				http.Error(w, "unable save file", http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fileSaved)

			return

		}

		defer f.Close()
		io.Copy(f, file)
	}
}
