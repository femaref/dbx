package dbx

import (
	"fmt"
	"reflect"
)

func Delete(target interface{}) error {
	db, err := Open()

	if err != nil {
		return err
	}
	return DeleteWithDB(db, target)
}

func DeleteWithDB(db DBAccess, target interface{}) error {
	assertPointerToStruct(target)

	t := reflect.TypeOf(target)
	t = t.Elem()

	id := getID(target)
	assertLiteral(id)

	stmt, _ := db.Preparex(fmt.Sprintf(deleteString, QuoteIdentifier(tableName(target, t))))
	defer stmt.Close()

	_, err := stmt.Exec(target, id)

	return err
}
