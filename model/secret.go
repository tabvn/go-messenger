package model

import (
	"messenger/db"
	"database/sql"
)

type Secret struct {
	Id     int64  `json:"id"`
	Secret string `json:"secret"`
}

func CheckSecret(secret string) (bool) {

	if secret == "" {
		return false
	}
	query := `SELECT COUNT(*) as c FROM secrets WHERE secret = ?`
	row, err := db.DB.FindOne(query, secret)

	if err != nil {
		return false
	}

	var count sql.NullInt64

	row.Scan(&count)

	if count.Int64 > 0 {
		return true
	}


	return false
}
