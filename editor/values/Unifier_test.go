package values_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/editor/values"

	"github.com/stretchr/testify/assert"
)

func TestUnifier_Integers(t *testing.T) {
	tt := []struct {
		values          []int
		expectedValue   interface{}
		expectedUnified bool
	}{
		{[]int{1, 1, 1}, 1, true},
		{[]int{}, nil, false},
		{[]int{1, 2, 3}, nil, false},
	}

	for tn, tc := range tt {
		u := values.NewUnifier()
		for _, value := range tc.values {
			u.Add(value)
		}
		resultValue := u.Unified()
		resultUnified := u.IsUnique()
		assert.Equal(t, tc.expectedValue, resultValue, fmt.Sprintf("Wrong result for test number %d", tn))
		assert.Equal(t, tc.expectedUnified, resultUnified, fmt.Sprintf("Wrong unified for test number %d", tn))
	}
}

func TestUnifier_Strings(t *testing.T) {
	tt := []struct {
		values          []string
		expectedValue   interface{}
		expectedUnified bool
	}{
		{[]string{"a", "a", "a"}, "a", true},
		{[]string{}, nil, false},
		{[]string{"a", "b", "c"}, nil, false},
	}

	for tn, tc := range tt {
		u := values.NewUnifier()
		for _, value := range tc.values {
			u.Add(value)
		}
		resultValue := u.Unified()
		resultUnified := u.IsUnique()
		assert.Equal(t, tc.expectedValue, resultValue, fmt.Sprintf("Wrong result for test number %d", tn))
		assert.Equal(t, tc.expectedUnified, resultUnified, fmt.Sprintf("Wrong unified for test number %d", tn))
	}
}

func TestUnifier_Types(t *testing.T) {
	type simpleType struct {
		value int
		text  string
	}
	tt := []struct {
		values          []simpleType
		expectedValue   interface{}
		expectedUnified bool
	}{
		{[]simpleType{{value: 1, text: "a"}, {value: 1, text: "a"}, {value: 1, text: "a"}}, simpleType{value: 1, text: "a"}, true},
		{[]simpleType{}, nil, false},
		{[]simpleType{{value: 1, text: "a"}, {value: 1, text: "a"}, {value: 2, text: "a"}}, nil, false},
		{[]simpleType{{value: 1, text: "a"}, {value: 1, text: "a"}, {value: 1, text: "b"}}, nil, false},
	}

	for tn, tc := range tt {
		u := values.NewUnifier()
		for _, value := range tc.values {
			u.Add(value)
		}
		resultValue := u.Unified()
		resultUnified := u.IsUnique()
		assert.Equal(t, tc.expectedValue, resultValue, fmt.Sprintf("Wrong result for test number %d", tn))
		assert.Equal(t, tc.expectedUnified, resultUnified, fmt.Sprintf("Wrong unified for test number %d", tn))
	}
}
