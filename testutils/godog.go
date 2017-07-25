package testutils

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"gopkg.in/mgo.v2/bson"
)

// TestMain -
func TestMain(m *testing.M, FeatureContext func(suite *godog.Suite)) {
	status := godog.RunWithOptions("godogs", FeatureContext, godog.Options{
		Format: "progress",
		Paths:  []string{"features"},
	})
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

// ParseTable receives a godog gherkin table and returns a map
// containing the rows of the table
func ParseTable(table *gherkin.DataTable) []map[string]string {
	if len(table.Rows) == 0 {
		return []map[string]string{}
	}

	headRow := table.Rows[0]

	valueRows := table.Rows[1:]
	values := make([]map[string]string, len(valueRows))
	for i := 0; i < len(valueRows); i++ {
		rowMap := map[string]string{}
		for i, cell := range valueRows[i].Cells {
			rowMap[headRow.Cells[i].Value] = cell.Value
		}
		values[i] = rowMap
	}

	return values
}

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
			fieldTime, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				return fmt.Errorf("Cannot set %s as a date", name)
			}

			structFieldValue.Set(reflect.ValueOf(fieldTime))
			continue

		case reflect.TypeOf([]string{}):
			structFieldValue.Set(reflect.ValueOf(strings.Split(value, ",")))
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
