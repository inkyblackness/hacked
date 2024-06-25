package ids

import "github.com/inkyblackness/hacked/ss1/resource"

// Palette identifier are listed below.
const (
	GamePalettesStart resource.ID = 0x02BC
)

// Texture identifier are listed below.
const (
	IconTextures   resource.ID = 0x004C
	SmallTextures  resource.ID = 0x004D
	MediumTextures resource.ID = 0x02C3
	LargeTextures  resource.ID = 0x03E8

	ScreenTextures resource.ID = 0x0141

	TextureNames  resource.ID = 0x086A
	TextureUsages resource.ID = 0x086B
)

// Bitmap identifier are listed below.
const (
	ObjectBitmaps         resource.ID = 0x0546
	ObjectTextureBitmaps  resource.ID = 0x01DB
	ObjectMaterialBitmaps resource.ID = 0x0884

	IconBitmaps     resource.ID = 0x004E
	GraffitiBitmaps resource.ID = 0x004F

	MfdDataBitmaps resource.ID = 0x0028
)

// Animation and video identifier are listed below.
const (
	VideoMailBitmapsStart    resource.ID = 0x0A40
	VideoMailAnimationsStart resource.ID = 0x0A4C
)

// Movie identifier are listed below.
const (
	MovieIntro resource.ID = 0x0BD6
	MovieDeath resource.ID = 0x0BD7
	MovieEnd   resource.ID = 0x0BD8
)

// Text identifier are listed below.
const (
	PaperTextsStart      resource.ID = 0x003C
	TrapMessageTexts     resource.ID = 0x0867
	WordTexts            resource.ID = 0x0868
	PanelNameTexts       resource.ID = 0x0869
	LogCategoryTexts     resource.ID = 0x0870
	VariousMessageTexts  resource.ID = 0x0871
	ScreenMessageTexts   resource.ID = 0x0877
	InfoNodeMessageTexts resource.ID = 0x0878
	AccessCardNameTexts  resource.ID = 0x0879
	DataletMessageTexts  resource.ID = 0x087A

	ObjectLongNames  resource.ID = 0x0024
	ObjectShortNames resource.ID = 0x086D
)

// Message identifier are listed below.
const (
	MailsStart     resource.ID = 0x0989
	LogsStart      resource.ID = 0x09B8
	FragmentsStart resource.ID = 0x0A98

	MailsAudioStart resource.ID = 0x0989 + 300
	LogsAudioStart  resource.ID = 0x09B8 + 300
)

// Sound identifier are listed below.
const (
	TrapMessagesAudioStart resource.ID = 0x0C1C

	SoundEffectsAudioStart resource.ID = 0x00C9
)

// Archive identifier are listed below.
const (
	ArchiveName resource.ID = 0x0FA0
	GameState   resource.ID = 0x0FA1

	LevelResourcesStart resource.ID = 4000
)
