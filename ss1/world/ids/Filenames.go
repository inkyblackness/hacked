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

// GameScr contains big bitmaps.
var GameScr = resource.AnyLanguage("gamescr.res")

// Texture contains all textures.
var Texture = resource.AnyLanguage("texture.res")

// VidMail contains all video mails.
var VidMail = resource.AnyLanguage("vidmail.res")

// Death contains another death video.
var Death = resource.AnyLanguage("death.res")

// Intro contains another intro video.
var Intro = resource.AnyLanguage("intro.res")

// Start1 contains ... ?.
var Start1 = resource.AnyLanguage("start1.res")

// Win1 contains ... ?.
var Win1 = resource.AnyLanguage("win1.res")

// LowDeth contains the low-res death video.
var LowDeth = resource.AnyLanguage("lowdeth.res")

// LowEnd contains the low-res end video.
var LowEnd = resource.AnyLanguage("lowend.res")

// SvgaDeth contains the high-res death video.
var SvgaDeth = resource.AnyLanguage("svgadeth.res")

// SvgaEnd contains the high-res end video.
var SvgaEnd = resource.AnyLanguage("svgaend.res")

// Obj3D contains 3D objects.
var Obj3D = resource.AnyLanguage("obj3D.res")

// ObjArt contains object art.
var ObjArt = resource.AnyLanguage("objart.res")

// ObjArt2 contains further object art.
var ObjArt2 = resource.AnyLanguage("objart2.res")

// ObjArt3 contains further object art.
var ObjArt3 = resource.AnyLanguage("objart3.res")

// CitMat contains materials for 3D objects.
var CitMat = resource.AnyLanguage("citmat.res")

// CutsPal contains palettes for the cutscenes.
var CutsPal = resource.AnyLanguage("cutspal.res")

// HandArt contains bitmaps for grabbed things.
var HandArt = resource.AnyLanguage("handart.res")

// SideArt contains bitmaps for the side buttons.
var SideArt = resource.AnyLanguage("sideart.res")

// DigiFX contains all the effect sounds.
var DigiFX = resource.AnyLanguage("digifx.res")

// Splash contains the splash screens.
var Splash = resource.AnyLanguage("splash.res")

// SplshPal contains the splash screen palettes.
var SplshPal = resource.AnyLanguage("splspal.res")

// LowResVideos returns the filename descriptors of all low-res videos.
func LowResVideos() resource.FilenameList {
	return []resource.Filename{LowIntr, LowDeth, LowEnd, Intro, Start1, Win1}
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
