package testutils

import (
	"container/list"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// FillStruct reads map 'm' and fills struct from pointer 's' with respective fields on 'm'
func FillStruct(s interface{}, m map[string]string) error {
	structValue := reflect.ValueOf(s).Elem()
	for name, value := range m {

		structFieldValue := structValue.FieldByName(name)
		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}

		switch structFieldValue.Type() {
		case reflect.TypeOf(bson.NewObjectId()):
			structFieldValue.Set(reflect.ValueOf(bson.ObjectIdHex(value)))
			continue

		case reflect.TypeOf(time.Time{}):
			var fieldTime time.Time
			var err error
			for _, layout := range supportedLayouts {
				fieldTime, err = time.Parse(layout, value)
				if err != nil {
					continue
				}

				break
			}

			if err != nil {
				return fmt.Errorf("cannot set %v as a date", name)
			}

			structFieldValue.Set(reflect.ValueOf(fieldTime))
			continue

		case reflect.TypeOf([]string{}):
			structFieldValue.Set(reflect.ValueOf(strings.Split(value, ",")))
			continue

		case reflect.TypeOf(0):
			fieldInt, err := strconv.Atoi(value)
			if err != nil {
				return err
			}

			structFieldValue.Set(reflect.ValueOf(fieldInt))
			continue

		default:
			structFieldValue.Set(reflect.ValueOf(value))
		}
	}
	return nil
}

// CreateSlice returns a new slice of type 't' filled with data from 'm' array of map
func CreateSlice(t interface{}, m []map[string]string) interface{} {
	kind := reflect.TypeOf(t)

	arr := reflect.MakeSlice(reflect.SliceOf(kind), 0, 0)

	for _, i := range m {
		v := reflect.New(kind)
		FillStruct(v.Interface(), i)
		arr = reflect.Append(arr, v.Elem())
	}
	return arr.Interface()
}

// JSONEqualsIgnoreOrder compares two jsons j1 and j2 ignoring their fields order
func JSONEqualsIgnoreOrder(j1, j2 interface{}) bool {
	if j1 == nil || j2 == nil {
		return j1 == j2
	}

	switch o1 := j1.(type) {
	case bool:
		return j1 == j2
	case string:
		return j1 == j2
	case float64:
		return j1 == j2
	case map[string]interface{}:
		o2, ok := j2.(map[string]interface{})
		if !ok {
			return false
		} else if len(o1) != len(o2) {
			return false
		}

		for k1, v1 := range o1 {
			if !JSONEqualsIgnoreOrder(v1, o2[k1]) {
				return false
			}
		}

		return true
	case []interface{}:
		o2, ok := j2.([]interface{})
		if !ok {
			return false
		} else if len(o1) != len(o2) {
			return false
		}

		diff := list.New()
		for _, v1 := range o1 {
			diff.PushBack(v1)
		}
		for _, v2 := range o2 {
			for e := diff.Front(); e != nil; e = e.Next() {
				if JSONEqualsIgnoreOrder(v2, e.Value) {
					diff.Remove(e)
					break
				}
			}
		}

		return diff.Len() == 0
	}

	panic("Unrecognized unmarshalled JSON type")
}

func sameStringsSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	xmap := make(map[string]int, len(x))
	for _, _x := range x {
		xmap[_x]++
	}
	ymap := make(map[string]int, len(y))
	for _, _y := range y {
		ymap[_y]++
	}
	return reflect.DeepEqual(xmap, ymap)
}
