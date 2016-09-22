package dbx

import (
	"errors"
	"fmt"
	"github.com/chuckpreslar/inflect"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

var (
    dialect string
	connectionString string
	isConfigured     bool

	NotConfiguredError error = errors.New("Database not configured")
)

func Configure(dialectx, connection string) {
    dialect = dialectx
	connectionString = connection
	isConfigured = true
}

func Open() (*sqlx.DB, error) {
	if !isConfigured {
		panic(NotConfiguredError)
	}

	return sqlx.Connect(dialect, connectionString)
}

type Tabler interface {
	TableName() string
}

func tableName(i interface{}, t reflect.Type) string {
	if tabler, ok := i.(Tabler); ok {
		return tabler.TableName()
	}

	plural := inflect.Pluralize(t.Name())
	return inflect.Underscore(plural)
}

func columnName(t reflect.StructField) string {
	tag := t.Tag

	value := tag.Get("db")

	if value != "" {
		return value
	}

	return inflect.Underscore(t.Name)
}

func generatePlaceholders(n int) string {
	var placeholders []string
	for i := 1; i <= n; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}

	return strings.Join(placeholders, ", ")
}

var insertString = `INSERT INTO %s (%s) VALUES (%s) RETURNING "id"`
var selectString = `SELECT * FROM %s WHERE "id" = $1`

var byteSlice reflect.Type = reflect.SliceOf(reflect.TypeOf(byte(0)))

func assertPointerToStruct(i interface{}) (bool) {
	t := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	} else {
		panic("argument not a pointer")
	}

	if t.Kind() != reflect.Struct {
		panic("argument target not a struct")
	}
	return true
}

func assertPointerToSlice(i interface{}) (bool) {
	t := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	} else {
		panic("argument not a pointer")
	}

	if t.Kind() != reflect.Slice {
		panic("argument target not a slice")
	}
	return true
}
const (
    int_literal = reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64
    uint_literal = reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64
    float_literal = reflect.Float64 | reflect.Float128
    literal reflect.Kind = reflect.Bool | reflect.String | int_literal | uint_literal | float_literal
)
func assertLiteral(i interface{}) (bool) {
    t := reflect.TypeOf(i)
    
    if t.Kind() & literal == 0 {
        panic("argument not a literal")
    }
    
    return true
}

func assertStruct(i interface{}) (bool) {
    t := reflect.TypeOf(i)
    
    if t.Kind() != reflect.Struct {
        panic("argument not a struct")
    }
    
    return true
}

func Create(i interface{}) (interface{}, error) {
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
		isNullValue := valField.Interface() == reflect.Zero(f.Type).Interface()
		isConvertibleToByteSlice := f.Type.ConvertibleTo(byteSlice)
		// only drop the id column and use it as target when it's the null value
		if column_name == "id" && isNullValue {
			id_val = val.Field(i)
			continue
		}
		
		columns = append(columns, column_name)
		
		if isConvertibleToByteSlice {
		    v := val.Field(i).Convert(reflect.TypeOf(""))
		    fields = append(fields, v.Interface())
		    continue
		}

		fields = append(fields, val.Field(i).Interface())
	}

	rendered_columns := strings.Join(columns, ", ")
	placeholders := generatePlaceholders(len(columns))

	prepared := fmt.Sprintf(insertString, table_name, rendered_columns, placeholders)

	db, err := Open()

	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(prepared)

	if err != nil {
		return nil, err
	}
	var v interface{}
	err = stmt.QueryRow(fields...).Scan(&v)

	if err != nil {
		return nil, err
	}

	if id_val.IsValid() {
		id_val.Set(reflect.ValueOf(v))
	}

	return v, nil
}

func Get(target interface{}, id interface{}) (error) {
    assertPointerToStruct(target)
    assertLiteral(id)
    
    t := reflect.TypeOf(i)
	t = t.Elem()
    
    return db.Get(target, fmt.Sprintf(selectString, tableName(i, t)), id)
}

func Select(target interface{}, cond interface{}) (error) {
    assertPointerToStruct(target)
    assertStruct(cond)
    
    tTarget := reflect.TypeOf(target)
	valTarget := reflect.ValueOf(target)
	
	tTarget = t.Elem()
	valTarget = val.Elem()
    
    tCond := reflect.TypeOf(cond)
	valCond := reflect.ValueOf(cond)
	
	if tTarget != tCond {
	    panic("target and cond type are different")
	}
	
	for i := 0; i < tCond.NumField(); i++ {
	    f := tCond.Field(i)
	    valField := valCond.Field(i)

		
		isNullValue := valField.Interface() == reflect.Zero(f.Type).Interface()
		if isNullValue {
		    continue
		}
		column_name := columnName(f)
		isConvertibleToByteSlice := f.Type.ConvertibleTo(byteSlice)
		
		columns = append(columns, column_name)
		
	    if isConvertibleToByteSlice {
		    v := val.Field(i).Convert(reflect.TypeOf(""))
		    fields = append(fields, v.Interface())
		    continue
		}

		fields = append(fields, val.Field(i).Interface())
	}    
}