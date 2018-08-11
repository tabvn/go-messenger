package helper

import "github.com/satori/go.uuid"

func GenerateID() (string) {

	return uuid.Must(uuid.NewV4()).String()
}
