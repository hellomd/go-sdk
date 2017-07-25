package testutils

import (
	"container/list"
	"reflect"
)

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
