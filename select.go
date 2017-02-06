package dbx

import (
	"errors"
	"reflect"
)

func Select(target interface{}, cond interface{}) error {
	assertPointerToStruct(target)
	assertStruct(cond)

	tTarget := reflect.TypeOf(target)
	valTarget := reflect.ValueOf(target)

	tTarget = tTarget.Elem()
	valTarget = valTarget.Elem()

	tCond := reflect.TypeOf(cond)
	valCond := reflect.ValueOf(cond)

	if tTarget != tCond {
		panic("target and cond type are different")
	}

	var columns []string
	var fields []interface{}

	for i := 0; i < tCond.NumField(); i++ {
		f := tCond.Field(i)
		valField := valCond.Field(i)

		mods := settings(f)
		if mods["ignore"] {
			continue
		}
		isNullValue := valField.Interface() == reflect.Zero(f.Type).Interface()
		if isNullValue {
			continue
		}
		column_name := columnName(f)
		isConvertibleToByteSlice := f.Type.ConvertibleTo(byteSlice)

		columns = append(columns, column_name)

		if isConvertibleToByteSlice {
			v := valCond.Field(i).Convert(reflect.TypeOf(""))
			fields = append(fields, v.Interface())
			continue
		}

		fields = append(fields, valCond.Field(i).Interface())
	}

	return errors.New("Not fully implemented")
}
