package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type Database struct {
	conn *sql.DB
}

var DB = &Database{

}

type MySQLConfig struct {
	Username   string
	Password   string
	Host       string
	Port       int
	UnixSocket string
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

func newDatabase(url string) (*Database, error) {

	conn, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	db := &Database{
		conn: conn,
	}

	DB = db

	return db, err
}

func (db *Database) Close() {
	DB.conn.Close()
}

func (db *Database) Query(query string, args interface{}) (*sql.Rows, error) {
	return DB.conn.Query(query, args)
}
func (db *Database) QueryRow(query string, args ...interface{}) (*sql.Row) {

	return DB.conn.QueryRow(query, args...)
}

func (db *Database) Prepare(query string) (*sql.Stmt, error) {
	return DB.conn.Prepare(query)
}

func InitDatabase(url string) (*Database, error) {

	DB, err := newDatabase(url)

	if err != nil {
		return nil, err
	}

	return DB, err
}

func (db *Database) Insert(query string, args ...interface{}) (int64, error) {

	/*stmt, prepareErr := DB.Prepare(query)

	if prepareErr != nil {
		return 0, prepareErr
	}

	defer stmt.Close()

	r, err := stmt.Exec(args...)

	*/


	r, err := DB.conn.Exec(query, args...)

	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	rowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return 0, fmt.Errorf("expected 1 row affected, got %d", rowsAffected)
	}

	lastInsertID, err := r.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertID, nil

}

func (db *Database) InsertMany(query string, args ...interface{}) (int64, error) {

	/*stmt, prepareErr := DB.Prepare(query)

	if prepareErr != nil {
		return 0, prepareErr
	}
	r, err := stmt.Exec(args...)

	*/

	r, err := DB.conn.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	rowsAffected, err := r.RowsAffected()

	return rowsAffected, nil

}
func (db *Database) Update(query string, args ...interface{}) (int64, error) {

	/*stmt, e := DB.conn.Prepare(query)
	if e != nil {
		return 0, e
	}
	//need close
	defer stmt.Close()

	r, err := stmt.Exec(args...)

	*/

	r, err := DB.conn.Exec(query, args...)

	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	rowsAffected, err := r.RowsAffected()

	return rowsAffected, err

}

func (db *Database) Count(query string, args ...interface{}) (int, error) {

	var count int

	row := DB.conn.QueryRow(query, args...)

	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) Get(table string, id int64) (*sql.Row, error) {

	query := "SELECT * FROM " +
		table +
		" WHERE id = ?"
	stmt, e := DB.conn.Prepare(query)
	if e != nil {
		return nil, e
	}

	row := stmt.QueryRow(id)

	return row, nil

}
func (db *Database) FindOne(query string, args ...interface{}) (*sql.Row, error) {

	/*stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(args...)
	*/

	row := DB.conn.QueryRow(query, args...)

	return row, nil
}

func (db *Database) Delete(query string, args ...interface{}) (int64, error) {

	/*stmt, e := DB.conn.Prepare(query)

	if e != nil {
		return 0, e
	}

	defer stmt.Close()

	r, err := stmt.Exec(query, args)
	*/

	r, err := DB.conn.Exec(query, args...)

	if err != nil {
		return 0, fmt.Errorf("could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()

	return rowsAffected, err

}

func (db *Database) DeleteMany(query string, args ...interface{}) (int64, error) {

	/*stmt, er := DB.conn.Prepare(query)

	if er != nil {
		return 0, er
	}

	defer stmt.Close()

	r, err := stmt.Exec(args...)
	*/

	r, err := DB.conn.Exec(query, args...)

	if err != nil {
		return 0, fmt.Errorf("could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("could not get rows affected: %v", err)
	}

	return rowsAffected, nil

}

func (db *Database) List(query string, args ...interface{}) (*sql.Rows, error) {

	/*stmt, err := DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)
	*/

	rows, err := DB.conn.Query(query, args...)
	return rows, err

}
