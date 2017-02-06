package dbx

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DBAccess interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	MustExec(string, ...interface{}) sql.Result
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error

	Prepare(query string) (*sql.Stmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
}
