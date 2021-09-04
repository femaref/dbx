package dbx

import (
	"fmt"
	"reflect"
)

func (db *DB) DeleteWithDB(target interface{}) error {
	assertPointerToStruct(target)

	t := reflect.TypeOf(target)
	t = t.Elem()

	id := getID(target)
	assertLiteral(id)

	stmt, _ := db.DB.Preparex(fmt.Sprintf(deleteString, db.QuoteIdentifier(tableName(target, t))))
	defer stmt.Close()

	_, err := stmt.Exec(target, id)

	return err
}
