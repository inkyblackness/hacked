package resource_test

import (
	"fmt"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLanguageIncludes(t *testing.T) {
	tt := []struct {
		primary   resource.Language
		secondary resource.Language
		expected  bool
	}{
		{resource.LangAny, resource.LangAny, true},
		{resource.LangAny, resource.LangGerman, true},
		{resource.LangFrench, resource.LangAny, false},
		{resource.LangGerman, resource.LangFrench, false},
		{resource.LangFrench, resource.LangFrench, true},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("Expecting %v including %v should be %v", tc.primary, tc.secondary, tc.expected), func(t *testing.T) {
			result := tc.primary.Includes(tc.secondary)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLanguageString(t *testing.T) {
	tt := []struct {
		lang     resource.Language
		expected string
	}{
		{resource.LangAny, "Any"},
		{resource.LangDefault, "Default"},
		{resource.LangFrench, "French"},
		{resource.LangGerman, "German"},
		{resource.Language(0x40), "Unknown40"},
	}

	for _, tc := range tt {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.lang.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLanguages(t *testing.T) {
	result := resource.Languages()
	assert.Equal(t, 3, len(result))
}
