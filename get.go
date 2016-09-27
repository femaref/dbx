package dbx

import (
	"fmt"
	"reflect"
)

func Get(target interface{}, id interface{}) (error) {
	db, err := Open()

	if err != nil {
		return err
	}
	defer db.Close()
	
	return GetWithDB(db, target, id)
}

func GetWithDB(db DBAccess, target interface{}, id interface{}) (error) {
    assertPointerToStruct(target)
    assertLiteral(id)
    
    t := reflect.TypeOf(target)
	t = t.Elem()
	

	
	stmt, _ := db.Preparex(fmt.Sprintf(selectString, tableName(target, t)))
    
    return stmt.Get(target, id)
}