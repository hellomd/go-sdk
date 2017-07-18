package steps

import (
	"fmt"

	"encoding/json"

	"reflect"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/hellomd/go-sdk/testutils"
	"github.com/hellomd/go-sdk/testutils/fakes"
)

// TheFollowingEventsShouldHaveBeenPublished checks whether the desired events have been published
//
// Usage:
// Then the following events should have been published
//   | Key                         | Body                                     |
//   | questions.article.published | {"id": "abc", "title": "First question"} |
func TheFollowingEventsShouldHaveBeenPublished(table *gherkin.DataTable, pub *fakes.Publisher) error {
	expected := testutils.ParseTable(table)
	actual := pub.GetPublished()

	expectedMaps := []map[string]interface{}{}
	for _, e := range expected {
		var body interface{}
		if err := json.Unmarshal([]byte(e["Body"]), &body); err != nil {
			return fmt.Errorf("failed to unmarshal expected body: %v", err)
		}

		eMap := map[string]interface{}{
			"Key":  e["Key"],
			"Body": body,
		}
		expectedMaps = append(expectedMaps, eMap)
	}

	actualMaps := []map[string]interface{}{}
	for _, e := range actual {
		var body interface{}
		if err := json.Unmarshal(e.Body, &body); err != nil {
			return fmt.Errorf("failed to unmarshal actual body: %v", err)
		}

		eMap := map[string]interface{}{
			"Key":  e.Key,
			"Body": body,
		}
		actualMaps = append(actualMaps, eMap)
	}

	if len(expectedMaps) != len(actualMaps) {
		return fmt.Errorf("expected %v events, but got %v: %v", len(expectedMaps), len(actualMaps), actualMaps)
	}

	if !reflect.DeepEqual(expectedMaps, actualMaps) {
		return fmt.Errorf("expected different events to have been published:\n  Expected %v\n  Actual %v", expectedMaps, actualMaps)
	}

	return nil
}
