package ids

import "github.com/inkyblackness/hacked/ss1/resource"

// Palettes
const (
	GamePalettesStart resource.ID = 0x02BC
)

// Textures
const (
	IconTextures   resource.ID = 0x004C
	SmallTextures  resource.ID = 0x004D
	MediumTextures resource.ID = 0x02C3
	LargeTextures  resource.ID = 0x03E8

	TextureNames  resource.ID = 0x086A
	TextureUsages resource.ID = 0x086B
)

// Bitmaps
const (
	ObjectBitmaps         resource.ID = 0x0546
	ObjectTextureBitmaps  resource.ID = 0x01DB
	ObjectMaterialBitmaps resource.ID = 0x0884

	MfdDataBitmaps resource.ID = 0x0028
)

// Texts
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

	ObjectLongNames resource.ID = 0x0024
)

// Messages
const (
	MailsStart     resource.ID = 0x0989
	LogsStart      resource.ID = 0x09B8
	FragmentsStart resource.ID = 0x0A98

	MailsAudioStart resource.ID = 0x0989 + 300
	LogsAudioStart  resource.ID = 0x09B8 + 300
)

// Sounds
const (
	TrapMessagesAudioStart resource.ID = 0x0C1C
)

// Archives
const (
	ArchiveName resource.ID = 0x0FA0
	GameState   resource.ID = 0x0FA1

	LevelResourcesStart resource.ID = 4000
)
