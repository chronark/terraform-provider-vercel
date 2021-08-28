package util

import (
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
	tests := []struct {
		input1   []string
		input2   []string
		expected []string
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "d"}, []string{"c"}},
		{[]string{"a", "b", "d"}, []string{"a", "b", "c"}, []string{"d"}},
		{[]string{"a", "b", "c"}, []string{"d", "e", "f"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{}, []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		actual := Difference(test.input1, test.input2)

		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("expected: %v, received: %v", test.expected, actual)
		}
	}
}
