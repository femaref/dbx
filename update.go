package dbx

import (
	"fmt"
	"reflect"
	"strings"
	"database/sql/driver"
)

func Update(target interface{}, newValues interface{}) (error) {
	db, err := Open()

	if err != nil {
		return err
	}
	defer db.Close()
	
	return UpdateWithDB(db, target, newValues)
}

func UpdateWithDB(db DBAccess, target interface{}, newValues interface{}) (error) {
    assertPointerToStruct(target)
    assertStruct(newValues)
    
    
    t := reflect.TypeOf(target)
	valTarget := reflect.ValueOf(target)
	
	t = t.Elem()
	valTarget = valTarget.Elem()
	
	valNew := reflect.ValueOf(newValues)
	
	table_name := tableName(target, t)
	var columns []string
	var fields []interface{}
	var id_val reflect.Value = reflect.ValueOf(nil)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		valTargetField := valTarget.Field(i)
		valNewField := valNew.Field(i)

		column_name := columnName(f)
		mods := settings(f)
		if mods["ignore"] {
		    continue
		}
		isNullValue := reflect.DeepEqual(valNewField.Interface(), reflect.Zero(f.Type).Interface())
		isConvertibleToByteSlice := f.Type.ConvertibleTo(byteSlice)
		implementsValuer := f.Type.Implements(valuer)
		// only drop the id column and use it as target when it's the null value
		if column_name == "id" {
			id_val = valTargetField
			continue
		}
		
		if implementsValuer {
		    valuer := valNewField.Interface().(driver.Valuer)
		    if x, _ := valuer.Value(); x == nil {
		        continue
		    }
		}
		
		if isNullValue {
		    continue
		}
		
		columns = append(columns, column_name)
		
		if isConvertibleToByteSlice {
		    v := valNewField.Convert(reflect.TypeOf(""))
		    fields = append(fields, v.Interface())
		    continue
		}

		fields = append(fields, valNewField.Interface())
	}
	
	if !id_val.IsValid() {
	    panic("Could not find id column in struct")
	}
	
	placeholders := generatePlaceholders(len(columns), 1)
	fields = append([]interface{}{id_val.Interface()}, fields...)
	
	var field_updaters []string
	
	for idx, column_name := range columns {
	    single := fmt.Sprintf("%s = %s", column_name, placeholders[idx])
	    field_updaters = append(field_updaters, single)
	}
	
	qs := fmt.Sprintf(updateString, table_name, strings.Join(field_updaters, ", "))
	stmt, err := db.Prepare(qs)

	if err != nil {
		return  err
	}
	
	_, err =  stmt.Exec(fields...)	
	
	if err != nil {
	    return err
	}
	
	return nil
}