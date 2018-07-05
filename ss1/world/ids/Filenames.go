package ids

import (
	"strings"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// CybStrng contains all strings.
var CybStrng = resource.I18nFile([resource.LanguageCount]string{"cybstrng.res", "frnstrng.res", "gerstrng.res"})

// MfdArt contains all MFD graphics.
var MfdArt = resource.I18nFile([resource.LanguageCount]string{"mfdart.res", "mfdfrn.res", "mfdger.res"})

// CitALog contains all log audio.
var CitALog = resource.I18nFile([resource.LanguageCount]string{"citalog.res", "frnalog.res", "geralog.res"})

// CitBark contains all bark audio.
var CitBark = resource.I18nFile([resource.LanguageCount]string{"citbark.res", "frnbark.res", "gerbark.res"})

// LowIntr contains the low-res intro video.
var LowIntr = resource.I18nFile([resource.LanguageCount]string{"lowintr.res", "lofrintr.res", "logeintr.res"})

// SvgaIntr contains the high-res intro video.
var SvgaIntr = resource.I18nFile([resource.LanguageCount]string{"svgaintr.res", "svfrintr.res", "svgeintr.res"})

// Archive contains the game world.
var Archive = resource.AnyLanguage("archive.dat")

// GamePal contains the game palettes.
var GamePal = resource.AnyLanguage("gamepal.res")

// Texture contains all textures.
var Texture = resource.AnyLanguage("texture.res")

// LowDeth contains the low-res death video.
var LowDeth = resource.AnyLanguage("lowdeth.res")

// LowEnd contains the low-res end video.
var LowEnd = resource.AnyLanguage("lowend.res")

// SvgaDeth contains the high-res death video.
var SvgaDeth = resource.AnyLanguage("svgadeth.res")

// SvgaEnd contains the high-res end video.
var SvgaEnd = resource.AnyLanguage("svgaend.res")

// LowResVideos returns the filename descriptors of all low-res videos.
func LowResVideos() resource.FilenameList {
	return []resource.Filename{LowIntr, LowDeth, LowEnd}
}

// HighResVideos returns the filename descriptors of all high-res videos.
func HighResVideos() resource.FilenameList {
	return []resource.Filename{SvgaIntr, SvgaDeth, SvgaEnd}
}

// LocalizedFiles returns the filename descriptors of all files that are localized.
func LocalizedFiles() []resource.Filename {
	return []resource.Filename{CybStrng, MfdArt, CitALog, CitBark, LowIntr, SvgaIntr}
}

// LocalizeFilename returns the language that the resource file would typically contain.
func LocalizeFilename(filename string) resource.Language {
	all := LocalizedFiles()
	lowercase := strings.ToLower(filename)
	result := resource.LangAny
	for _, lang := range resource.Languages() {
		for _, file := range all {
			if file.For(lang) == lowercase {
				result = lang
			}
		}
	}
	return result
}
