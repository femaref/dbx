package dbx

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

var insertString = `INSERT INTO %s (%s) VALUES (%s)`
var insertStringWithID = `INSERT INTO %s (%s) VALUES (%s) RETURNING "id"`

func (db *DB) CreateWithDB(i interface{}) error {
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

		/*
			is null, nullable, wants null -> skipIfNull == true
			is null, nullable, wants default -> skipIfNull == false
			is null, expects non null -> db exception
		*/
		if isNullValue && mods["skipIfNull"] {
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

	for i := range columns {
		columns[i] = db.QuoteIdentifier(columns[i])
	}
	rendered_columns := strings.Join(columns, ", ")
	placeholders := strings.Join(generatePlaceholders(len(columns), 0), ", ")

	insert_query := insertString
	if id_val.IsValid() {
		insert_query = insertStringWithID
	}
	prepared := fmt.Sprintf(insert_query, db.QuoteIdentifier(table_name), rendered_columns, placeholders)

	stmt, err := db.DB.Preparex(prepared)

	if err != nil {
		return err
	}
	defer stmt.Close()

	var v interface{}
	if id_val.IsValid() {
		err = stmt.QueryRow(fields...).Scan(&v)
		if err != nil {
			return err
		}
		id_val.Set(reflect.ValueOf(v))
	} else {
		_, err = stmt.Exec(fields...)
		if err != nil {
			return err
		}
	}

	return nil
}
