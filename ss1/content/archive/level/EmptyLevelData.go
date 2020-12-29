package level

import (
	"bytes"
	"math/bits"

	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// EmptyLevelParameters contain the values to create an empty level.
type EmptyLevelParameters struct {
	// XSize specifies the expected width, in tiles. 0 (or less) will default.
	XSize int32
	// YSize specifies the expected height, in tiles. 0 (or less) will default.
	YSize int32
	// Cyberspace indicates whether the level should be marked as cyberspace.
	Cyberspace bool
	// MapModifier is called to make initial changes to the map before serializing.
	// This can be used to empty out starter tile on starter level.
	MapModifier func(x, y int, entry *TileMapEntry)
}

// EmptyLevelData returns an array of serialized data for an empty level.
func EmptyLevelData(param EmptyLevelParameters) [lvlids.PerLevel][]byte {
	var levelData [lvlids.PerLevel][]byte
	baseInfo := DefaultBaseInfo(param.Cyberspace)
	if param.XSize > 0 {
		baseInfo.XSize = param.XSize
		baseInfo.XShift = int32(bits.Len32(uint32(param.XSize)))
		if (1 << (baseInfo.XShift - 1)) == baseInfo.XSize {
			baseInfo.XShift--
		}
	}
	if param.YSize > 0 {
		baseInfo.YSize = param.YSize
		baseInfo.YShift = int32(bits.Len32(uint32(param.YSize)))
		if (1 << (baseInfo.YShift - 1)) == baseInfo.YSize {
			baseInfo.YShift--
		}
	}

	levelData[lvlids.MapVersionNumber] = encode(mapVersionValue)
	levelData[lvlids.ObjectVersionNumber] = encode(objectVersionValue)
	levelData[lvlids.Information] = encode(&baseInfo)

	tileMap := NewTileMap(1<<baseInfo.XShift, 1<<baseInfo.YShift)
	if param.MapModifier != nil {
		for y := 0; y < int(baseInfo.YSize); y++ {
			for x := 0; x < int(baseInfo.XSize); x++ {
				param.MapModifier(x, y, tileMap.Tile(x, y, int(baseInfo.XShift)))
			}
		}
	}
	levelData[lvlids.TileMap] = encode(tileMap)

	levelData[lvlids.Schedules] = encode(make([]byte, baseInfo.Scheduler.ElementSize*1))
	levelData[lvlids.TextureAtlas] = encode(make(TextureAtlas, DefaultTextureAtlasSize))
	levelData[lvlids.ObjectMainTable] = encode(DefaultObjectMainTable())
	levelData[lvlids.ObjectCrossRefTable] = encode(DefaultObjectCrossReferenceTable())
	for class := object.Class(0); class < object.ClassCount; class++ {
		levelData[lvlids.ObjectClassTablesStart+int(class)] = encode(DefaultObjectClassTable(class))
		levelData[lvlids.ObjectDefaultTablesStart+int(class)] = encode(NewObjectClassEntry(ObjectClassInfoFor(class).DataSize))
	}

	levelData[lvlids.SavefileVersion] = encode(savefileVersionValue)
	levelData[lvlids.Unused41] = encode(UnknownL41{})

	levelData[lvlids.TextureAnimations] = encode([TextureAnimationCount]TextureAnimationEntry{})
	levelData[lvlids.SurveillanceSources] = encode(make([]ObjectID, SurveillanceObjectCount))
	levelData[lvlids.SurveillanceSurrogates] = encode(make([]ObjectID, SurveillanceObjectCount))
	levelData[lvlids.Parameters] = encode(DefaultParameters())
	levelData[lvlids.MapNotes] = encode(make([]byte, MapNotesSize))
	levelData[lvlids.MapNotesPointer] = encode(MapNotesPointer(0))

	levelData[lvlids.Unknown48] = encode(UnknownL48{})
	levelData[lvlids.Unknown49] = encode(UnknownL49{})
	levelData[lvlids.Unknown50] = encode(UnknownL50{})

	levelData[lvlids.LoopConfiguration] = encode(make([]byte, LoopConfigEntrySize*LoopConfigEntryCount))

	levelData[lvlids.Unknown52] = encode(UnknownL52{})
	levelData[lvlids.HeightSemaphores] = encode(HeightSemaphores{})

	return levelData
}

func encode(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(data)
	return buf.Bytes()
}
