package resource_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/resource"
)

func TestContentTypeString(t *testing.T) {
	tt := []struct {
		contentType resource.ContentType
		expected    string
	}{
		{resource.Palette, "Palette"},
		{resource.Text, "Text"},
		{resource.Bitmap, "Bitmap"},
		{resource.Font, "Font"},
		{resource.Animation, "Animation"},
		{resource.Sound, "Sound"},
		{resource.Geometry, "Geometry"},
		{resource.Movie, "Movie"},
		{resource.Archive, "Archive"},
		{resource.ContentType(254), "UnknownFE"},
	}

	for _, tc := range tt {
		result := tc.contentType.String()
		assert.Equal(t, tc.expected, result, fmt.Sprintf("Failed for 0x%02X", int(tc.contentType)))
	}
}
