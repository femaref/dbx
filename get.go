package dbx

import (
	"fmt"
	"reflect"
)

func (db *DB) Get(target interface{}, id interface{}) error {
	assertPointerToStruct(target)
	assertLiteral(id)

	t := reflect.TypeOf(target)
	t = t.Elem()

	stmt, _ := db.DB.Preparex(fmt.Sprintf(selectString, db.QuoteIdentifier(tableName(target, t))))
	defer stmt.Close()

	return stmt.Get(target, id)
}
