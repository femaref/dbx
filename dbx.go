package dbx

import (
	"errors"
	"fmt"
	"database/sql/driver"
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

func settings(t reflect.StructField) map[string]bool {
	tag := t.Tag

	value := tag.Get("dbx")
	
	var output map[string]bool = make(map[string]bool)

	if value != "" {
		split := strings.Split(value, ",")
		for _, cur := range split {
		    cur = strings.TrimSpace(cur)
		    
		    output[cur] = true
		}
	}

	return output
}

func generatePlaceholders(n int, offset int) []string {
	var placeholders []string
	for i := 1 + offset; i <= n + offset; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}

	return placeholders
}

var selectString = `SELECT * FROM %s WHERE "id" = $1`
var updateString = `UPDATE %s SET %s WHERE "id" = $1`

var byteSlice reflect.Type = reflect.SliceOf(reflect.TypeOf(byte(0)))
var valuer reflect.Type = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

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
var (
    int_literal = reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64
    uint_literal = reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64
    float_literal = reflect.Float32 | reflect.Float64
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


