#pragma once

#include <stdbool.h>

#include "hacked/core/media/Bitmap.h"
#include "hacked/core/media/Text.h"
#include "hacked/core/system/Validation.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct
{
   Codepoint codepoint;
   PixelRect rect;
} FontCodepointEntry;

typedef struct Font Font;
typedef void (*FontReleaseFunc)(Font const *font);

typedef struct Font
{
   Bitmap atlas;
   bool color;
   PixelAxisSize height;

   size_t codepointCount;
   /**
    * A list of codepoints supported by this font, with their pixel location values into the atlas.
    * The list must be ordered by the codepoint values, to allow for binary search.
    */
   FontCodepointEntry *codepoints;

   void *userData;
   /**
    * release function for the font, called from fontRelease().
    * This function is expected to release the userdata, the atlas, and the corresponding font structure itself.
    */
   FontReleaseFunc release;
} Font;

/**
 * Releases the font instance. In case it was a dynamically allocated memory, it will be freed.
 *
 * @param fontRef address of the font pointer to release; the pointed-to-pointer will be set to @code NULL@endcode.
 */
extern void fontRelease(Font const **fontRef);

extern ValidationResult fontValidate(Font const *font);

extern void fontSortCodepoints(FontCodepointEntry *entries, size_t count);

/**
 * Finds the entry of the font for given codepoint.
 *
 * @param font the font to query
 * @param codepoint the codepoint to find
 * @return the entry, if existing, or NULL
 */
extern FontCodepointEntry const *fontFindCodepointEntry(Font const *font, Codepoint codepoint);

#ifdef __cplusplus
}
#endif
