package helper

import (
	"reflect"
)

func GetIds(input interface{}) ([] int64) {

	var arrIds [] int64

	switch reflect.TypeOf(input).Kind() {

	case reflect.Slice:
		s := reflect.ValueOf(input)

		for i := 0; i < s.Len(); i++ {
			value := s.Index(i).Interface()

			arrIds = append(arrIds, int64(value.(int)))
		}
	}

	return arrIds
}
