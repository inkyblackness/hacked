package resource_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestLocalizedResources(t *testing.T) {
	tt := []struct {
		filename string
		expected resource.Language
	}{
		{"citmat.res", resource.LangAny},
		{"digifx.res", resource.LangAny},
		{"gamescr.res", resource.LangAny},
		{"handart.res", resource.LangAny},
		{"intro.res", resource.LangAny},
		{"obj3d.res", resource.LangAny},
		{"objart2.res", resource.LangAny},
		{"objart3.res", resource.LangAny},
		{"sideart.res", resource.LangAny},
		{"texture.res", resource.LangAny},

		{"cybstrng.res", resource.LangDefault},
		{"mfdart.res", resource.LangDefault},

		{"frnstrng.res", resource.LangFrench},
		{"mfdfrn.res", resource.LangFrench},

		{"gerstrng.res", resource.LangGerman},
		{"mfdger.res", resource.LangGerman},

		{"archive.dat", resource.LangAny},
		{"cutspal.res", resource.LangAny},
		{"death.res", resource.LangAny},
		{"gamepal.res", resource.LangAny},
		{"intro.res", resource.LangAny},
		{"lowdeth.res", resource.LangAny},
		{"lowend.res", resource.LangAny},
		{"objart.res", resource.LangAny},
		{"splash.res", resource.LangAny},
		{"splshpal.res", resource.LangAny},
		{"start1.res", resource.LangAny},
		{"svgadeth.res", resource.LangAny},
		{"svgaend.res", resource.LangAny},
		{"vidmail.res", resource.LangAny},
		{"win1.res", resource.LangAny},

		{"citalog.res", resource.LangDefault},
		{"citbark.res", resource.LangDefault},
		{"lowintr.res", resource.LangDefault},
		{"svgaintr.res", resource.LangDefault},

		{"frnalog.res", resource.LangFrench},
		{"frnbark.res", resource.LangFrench},
		{"lofrintr.res", resource.LangFrench},
		{"svfrintr.res", resource.LangFrench},

		{"geralog.res", resource.LangGerman},
		{"gerbark.res", resource.LangGerman},
		{"logeintr.res", resource.LangGerman},
		{"svgeintr.res", resource.LangGerman},
	}

	for _, tc := range tt {
		result := resource.LocalizeResourcesByFilename(nil, tc.filename)
		assert.Equal(t, tc.expected, result.Language, "Wrong language for <"+tc.filename+">")
	}
}
