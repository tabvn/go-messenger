package model

import (
	"database/sql"
	"errors"
	"messenger/db"
	"messenger/helper"
)

type File struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Name     string `json:"name"`
	Original string `json:"original"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Created  int64  `json:"created"`
}

func scanFiles(rows *sql.Rows) ([] *File, error) {

	var files [] *File

	var (
		id       sql.NullInt64
		userId   sql.NullInt64
		name     sql.NullString
		original sql.NullString
		fileType sql.NullString
		size     sql.NullInt64
		created  sql.NullInt64
	)


	for rows.Next() {

		err := rows.Scan(&id, &userId, &name, &original, &fileType, &size, &created)
		if err != nil {


		}

		file := &File{
			Id:       id.Int64,
			UserId:   userId.Int64,
			Name:     name.String,
			Original: original.String,
			Type:     fileType.String,
			Size:     size.Int64,
			Created:  created.Int64,
		}

		files = append(files, file)
	}

	return files, nil

}

func SaveFile(userId int64, name string, original string, fileType string, size int64) (*File, error) {

	created := helper.GetUnixTimestamp()

	q := `INSERT INTO files (user_id, name, original, type, size, created) VALUES (?, ?, ?, ?, ?, ?)`
	id, err := db.DB.Insert(q, userId, name, original, fileType, size, created)

	if err != nil {
		return nil, err
	}
	if id > 0 {
		return &File{
			Id:       id,
			UserId:   userId,
			Name:     name,
			Original: original,
			Type:     fileType,
			Size:     size,
			Created:  created,
		}, nil
	}
	return nil, errors.New("unknown error")
}
