package dbx

import (
	"fmt"
	"reflect"
	"strings"
	"database/sql/driver"
)


var insertString = `INSERT INTO %s (%s) VALUES (%s)`
var insertStringWithID = `INSERT INTO %s (%s) VALUES (%s) RETURNING "id"`
func Create(i interface{}) (interface{}, error) {
	db, err := Open()

	if err != nil {
		return nil, err
	}
	defer db.Close()
	
	return CreateWithDB(db, i)
}

func CreateWithDB(db DBAccess, i interface{}) (interface{}, error) {
	assertPointerToStruct(i)
	
	t := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	
	t = t.Elem()
	val = val.Elem()
	

	table_name := tableName(i, t)
	var columns []string
	var fields []interface{}
	var id_val reflect.Value = reflect.ValueOf(nil)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		valField := val.Field(i)

		column_name := columnName(f)
		mods := settings(f)
		if mods["ignore"] {
		    continue
		}
		isNullValue := reflect.DeepEqual(valField.Interface(), reflect.Zero(f.Type).Interface())
		isConvertibleToByteSlice := f.Type.ConvertibleTo(byteSlice)
		implementsValuer := f.Type.Implements(valuer)
		// only drop the id column and use it as target when it's the null value
		if column_name == "id" && isNullValue {
			id_val = val.Field(i)
			continue
		}
		if implementsValuer {
		    valuer := valField.Interface().(driver.Valuer)
		    x, _ := valuer.Value()
		    if x == nil {
		        continue
		    }
		}
		
		columns = append(columns, column_name)
		if !implementsValuer && isConvertibleToByteSlice {
		    v := val.Field(i).Convert(reflect.TypeOf(""))
		    fields = append(fields, v.Interface())
		    continue
		}
		
		fields = append(fields, val.Field(i).Interface())
	}

    for i, _ := range columns {
        columns[i] = QuoteIdentifier(columns[i])
    }
	rendered_columns := strings.Join(columns, ", ")
	placeholders := strings.Join(generatePlaceholders(len(columns), 0), ", ")

    insert_query := insertString
    if id_val.IsValid() {
        insert_query = insertStringWithID
    }
    prepared := fmt.Sprintf(insert_query, QuoteIdentifier(table_name), rendered_columns, placeholders)

	stmt, err := db.Preparex(prepared)

	if err != nil {
		return nil, err
	}

  	var v interface{}
	if id_val.IsValid() {
    	err = stmt.QueryRow(fields...).Scan(&v)
        if err != nil {
            return nil, err
        }
		id_val.Set(reflect.ValueOf(v))
	} else {
	    _, err = stmt.Exec(fields...)
	    if err != nil {
            return nil, err
        }
	}

	return v, nil
}