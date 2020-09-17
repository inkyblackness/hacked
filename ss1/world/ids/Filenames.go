package ids

import (
	"strings"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// CybStrng contains all strings.
var CybStrng = I18nFile([resource.LanguageCount]string{"cybstrng.res", "frnstrng.res", "gerstrng.res"})

// MfdArt contains all MFD graphics.
var MfdArt = I18nFile([resource.LanguageCount]string{"mfdart.res", "mfdfrn.res", "mfdger.res"})

// CitALog contains all log audio.
var CitALog = I18nFile([resource.LanguageCount]string{"citalog.res", "frnalog.res", "geralog.res"})

// CitBark contains all bark audio.
var CitBark = I18nFile([resource.LanguageCount]string{"citbark.res", "frnbark.res", "gerbark.res"})

// LowIntr contains the low-res intro video.
var LowIntr = I18nFile([resource.LanguageCount]string{"lowintr.res", "lofrintr.res", "logeintr.res"})

// SvgaIntr contains the high-res intro video.
var SvgaIntr = I18nFile([resource.LanguageCount]string{"svgaintr.res", "svfrintr.res", "svgeintr.res"})

// Archive contains the game world.
var Archive = AnyLanguage("archive.dat")

// GamePal contains the game palettes.
var GamePal = AnyLanguage("gamepal.res")

// GameScr contains big bitmaps.
var GameScr = AnyLanguage("gamescr.res")

// Texture contains all textures.
var Texture = AnyLanguage("texture.res")

// VidMail contains all video mails.
var VidMail = AnyLanguage("vidmail.res")

// Death contains another death video.
var Death = AnyLanguage("death.res")

// Intro contains another intro video.
var Intro = AnyLanguage("intro.res")

// Start1 contains ... ?.
var Start1 = AnyLanguage("start1.res")

// Win1 contains ... ?.
var Win1 = AnyLanguage("win1.res")

// LowDeth contains the low-res death video.
var LowDeth = AnyLanguage("lowdeth.res")

// LowEnd contains the low-res end video.
var LowEnd = AnyLanguage("lowend.res")

// SvgaDeth contains the high-res death video.
var SvgaDeth = AnyLanguage("svgadeth.res")

// SvgaEnd contains the high-res end video.
var SvgaEnd = AnyLanguage("svgaend.res")

// Obj3D contains 3D objects.
var Obj3D = AnyLanguage("obj3D.res")

// ObjArt contains object art.
var ObjArt = AnyLanguage("objart.res")

// ObjArt2 contains further object art.
var ObjArt2 = AnyLanguage("objart2.res")

// ObjArt3 contains further object art.
var ObjArt3 = AnyLanguage("objart3.res")

// CitMat contains materials for 3D objects.
var CitMat = AnyLanguage("citmat.res")

// CutsPal contains palettes for the cutscenes.
var CutsPal = AnyLanguage("cutspal.res")

// HandArt contains bitmaps for grabbed things.
var HandArt = AnyLanguage("handart.res")

// SideArt contains bitmaps for the side buttons.
var SideArt = AnyLanguage("sideart.res")

// DigiFX contains all the effect sounds.
var DigiFX = AnyLanguage("digifx.res")

// Splash contains the splash screens.
var Splash = AnyLanguage("splash.res")

// SplshPal contains the splash screen palettes.
var SplshPal = AnyLanguage("splspal.res")

// LowResVideos returns the filename descriptors of all low-res videos.
func LowResVideos() FilenameList {
	return []Filename{LowIntr, LowDeth, LowEnd, Intro, Start1, Win1}
}

// HighResVideos returns the filename descriptors of all high-res videos.
func HighResVideos() FilenameList {
	return []Filename{SvgaIntr, SvgaDeth, SvgaEnd}
}

// LocalizedFiles returns the filename descriptors of all files that are localized.
func LocalizedFiles() []Filename {
	return []Filename{CybStrng, MfdArt, CitALog, CitBark, LowIntr, SvgaIntr}
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
