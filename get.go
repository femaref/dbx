package dbx

import (
	"fmt"
	"reflect"
)

func Get(target interface{}, id interface{}) error {
	db, err := Open()

	if err != nil {
		return err
	}
	return GetWithDB(db, target, id)
}

func GetWithDB(db DBAccess, target interface{}, id interface{}) error {
	assertPointerToStruct(target)
	assertLiteral(id)

	t := reflect.TypeOf(target)
	t = t.Elem()

	stmt, _ := db.Preparex(fmt.Sprintf(selectString, QuoteIdentifier(tableName(target, t))))
	defer stmt.Close()

	return stmt.Get(target, id)
}
