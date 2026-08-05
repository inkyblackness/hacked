#pragma once

#include <stdbool.h>

#include "hacked/core/media/Bitmap.h"
#include "hacked/core/media/Text.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct
{
   Codepoint codepoint;
   PixelPosition topLeft;
   PixelSize size;
} FontCodepointEntry;

typedef struct Font Font;
typedef void (*FontReleaseFunc)(Font *font);

typedef struct Font
{
   Bitmap bitmap;
   bool color;

   size_t codepointCount;
   /**
    * A list of codepoints supported by this font, with their pixel location values into the bitmap.
    * The list must be ordered by the codepoint values, to allow for binary search.
    */
   FontCodepointEntry *codepoints;

   void *userData;
   /**
    * release function for the font, called from fontRelease().
    * This function is expected to release the userdata, and the corresponding font structure itself.
    */
   FontReleaseFunc release;
} Font;

/**
 * Releases the font instance. In case it was a dynamically allocated memory, it will be freed.
 *
 * @param fontRef address of the font pointer to release; the pointed-to-pointer will be set to @code NULL@endcode.
 */
extern void fontRelease(Font **fontRef);

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
