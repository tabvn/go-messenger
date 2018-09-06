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

			t := reflect.TypeOf(value)
			switch t.String() {

			case "float64":
				arrIds = append(arrIds, int64(value.(float64)))
				break

			case "int":

				arrIds = append(arrIds, int64(value.(int)))

				break

			default:
				break
			}

		}
	}

	return arrIds
}