#include <stdio.h>

#include "hacked/core/editor/ui/UIFont.h"
#include "hacked/core/media/PixelSpace.h"

#define NK_IMPLEMENTATION
#include "hacked/nuklear/NuklearSdlBridge.h"

struct UIBridgeDevice
{
   struct nk_buffer cmds;
   struct nk_draw_null_texture nullTexture;
   SDL_Texture *fontTexture;
};

struct UIBridgeVertex
{
   float position[2];
   float uv[2];
   float col[4];
};

struct UIBridge
{
   SDL_Window *window;
   SDL_Renderer *renderer;
   Font const *uiFont;
   struct nk_user_font *nkFont;
   struct UIBridgeDevice device;
   struct nk_context ctx;
   struct nk_allocator allocator;
   nk_handle userdata;
   Uint64 lastRenderTick;
   bool insertKeyToggle;
   bool editWasActive;
};

static void uiBridgeDeviceDropFontTexture(struct UIBridgeDevice *const device)
{
   if (device->fontTexture != NULL)
   {
      SDL_DestroyTexture(device->fontTexture);
      device->fontTexture = NULL;
   }
}

nk_handle uiBridgeGetUserdata(struct nk_context const *const ctx)
{
   NK_ASSERT(ctx);
   struct UIBridge *bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);
   return bridge->userdata;
}

void uiBridgeSetUserdata(struct nk_context *const ctx, nk_handle const userdata)
{
   NK_ASSERT(ctx);
   struct UIBridge *const bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);
   bridge->userdata = userdata;
}

static void *nkAllocatorSDLAlloc(nk_handle const user, void *const old, nk_size const size)
{
   NK_UNUSED(user);
   /* FIXME: should use SDL_realloc here, not SDL_malloc
    * but this could cause a double-free due to bug within Nuklear, see:
    * https://github.com/Immediate-Mode-UI/Nuklear/issues/768
    * */
#if 0
    return SDL_realloc(old, size);
#else
   NK_UNUSED(old);
   return SDL_malloc(size);
#endif
}

static void nkAllocatorSDLFree(nk_handle const user, void *const old)
{
   NK_UNUSED(user);
   SDL_free(old);
}

static struct nk_allocator nkAllocatorFromSDL()
{
   struct nk_allocator allocator = {0};
   allocator.userdata.ptr = NULL;
   allocator.alloc = nkAllocatorSDLAlloc;
   allocator.free = nkAllocatorSDLFree;
   return allocator;
}

static void uiBridgeUploadAtlas(struct nk_context *const ctx, void const *const image, int const width, int const height)
{
   NK_ASSERT(ctx);
   NK_ASSERT(image);
   NK_ASSERT(width > 0);
   NK_ASSERT(height > 0);

   struct UIBridge *bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);

   uiBridgeDeviceDropFontTexture(&bridge->device);

   bridge->device.fontTexture = SDL_CreateTexture(bridge->renderer, SDL_PIXELFORMAT_ARGB8888, SDL_TEXTUREACCESS_STATIC, width, height);
   NK_ASSERT(bridge->device.fontTexture);
   SDL_UpdateTexture(bridge->device.fontTexture, NULL, image, 4 * width);
   SDL_SetTextureBlendMode(bridge->device.fontTexture, SDL_BLENDMODE_BLEND);
   SDL_SetTextureScaleMode(bridge->device.fontTexture, SDL_SCALEMODE_NEAREST);
}

void uiBridgeUpdateTextInput(struct nk_context *ctx)
{
   bool active = false;
   NK_ASSERT(ctx);
   struct UIBridge *const bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);

   /* Determine if Nuklear is using any top-level "edit" widget.
    * Popups take higher priority because they block any incoming input.
    * This will not work, if the widget is not updating context state properly. */
   if (!ctx->active)
      active = false;
   else if (ctx->active->popup.win)
      active = ctx->active->popup.win->edit.active;
   else
      active = ctx->active->edit.active;

   /* decide, if TextInputActive should be unchanged/stoped/started
    * and change its state accordingly for owned SDL Window */
   if (active != bridge->editWasActive)
   {
      bool const window_edit_active = SDL_TextInputActive(bridge->window);

      /* If you ever hit this check, it means that the demo and your app
       * (or something else) are all trying to manage TextInputActive state.
       * This can cause subtle bugs where the state won't be what you expect.
       * You can safely remove this assert and the demo will keep working,
       * but make sure it does not cause any issues for you */
      NK_ASSERT(window_edit_active == bridge->editWasActive && "something else changed TextInputActive state for this Window");

      if (!window_edit_active && !bridge->editWasActive && active)
         SDL_StartTextInput(bridge->window);
      else if (window_edit_active && bridge->editWasActive && !active)
         SDL_StopTextInput(bridge->window);
      bridge->editWasActive = active;
   }

   /* FIXME:
    * for full SDL3 integration, you also need to find current edit widget
    * bounds and the text cursor offset, and pass this data into SDL_SetTextInputArea.
    * This is currently not possible to do safely as Nuklear does not support it.
    * https://wiki.libsdl.org/SDL3/SDL_SetTextInputArea
    * https://github.com/Immediate-Mode-UI/Nuklear/pull/857
    */
}

void uiBridgeRender(struct nk_context *const ctx)
{
   NK_ASSERT(ctx);
   struct UIBridge *bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);

   { /* setup internal delta time that Nuklear needs for animations */
      Uint64 const now = SDL_GetTicksNS();
      ctx->delta_time_seconds = (float)(now - bridge->lastRenderTick) / (float)SDL_NS_PER_SECOND;
      bridge->lastRenderTick = now;
   }

   {
      static enum nk_anti_aliasing const antiAliasing = NK_ANTI_ALIASING_ON;
      int const vs = sizeof(struct UIBridgeVertex);
      size_t const vp = NK_OFFSETOF(struct UIBridgeVertex, position);
      size_t const vt = NK_OFFSETOF(struct UIBridgeVertex, uv);
      size_t const vc = NK_OFFSETOF(struct UIBridgeVertex, col);

      /* convert from command queue into draw list and draw to screen */

      /* fill converting configuration */
      static const struct nk_draw_vertex_layout_element vertex_layout[] = {{NK_VERTEX_POSITION, NK_FORMAT_FLOAT, NK_OFFSETOF(struct UIBridgeVertex, position)},
         {NK_VERTEX_TEXCOORD, NK_FORMAT_FLOAT, NK_OFFSETOF(struct UIBridgeVertex, uv)},
         {NK_VERTEX_COLOR, NK_FORMAT_R32G32B32A32_FLOAT, NK_OFFSETOF(struct UIBridgeVertex, col)}, {NK_VERTEX_LAYOUT_END}};
      struct nk_convert_config config;
      NK_MEMSET(&config, 0, sizeof(config));
      config.vertex_layout = vertex_layout;
      config.vertex_size = sizeof(struct UIBridgeVertex);
      config.vertex_alignment = NK_ALIGNOF(struct UIBridgeVertex);
      config.tex_null = bridge->device.nullTexture;
      config.circle_segment_count = 22;
      config.curve_segment_count = 22;
      config.arc_segment_count = 22;
      config.global_alpha = 1.0f;
      config.shape_AA = antiAliasing;
      config.line_AA = antiAliasing;

      /* convert shapes into vertices */
      struct nk_buffer vbuf, ebuf;
      nk_buffer_init(&vbuf, &bridge->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
      nk_buffer_init(&ebuf, &bridge->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
      nk_convert(&bridge->ctx, &bridge->device.cmds, &vbuf, &ebuf, &config);

      /* iterate over and execute each draw command */
      nk_draw_index const *offset = (const nk_draw_index *)nk_buffer_memory_const(&ebuf);

      bool clipping_enabled = SDL_RenderClipEnabled(bridge->renderer);
      SDL_Rect saved_clip;
      SDL_GetRenderClipRect(bridge->renderer, &saved_clip);

      struct nk_draw_command const *cmd = NULL;
      nk_draw_foreach(cmd, &bridge->ctx, &bridge->device.cmds)
      {
         if (!cmd->elem_count)
            continue;

         {
            SDL_Rect r;
            r.x = (int)nk_roundf(cmd->clip_rect.x);
            r.y = (int)nk_roundf(cmd->clip_rect.y);
            r.w = (int)nk_roundf(cmd->clip_rect.w);
            r.h = (int)nk_roundf(cmd->clip_rect.h);
            SDL_SetRenderClipRect(bridge->renderer, &r);
         }

         {
            const void *vertices = nk_buffer_memory_const(&vbuf);

            SDL_RenderGeometryRaw(bridge->renderer, (SDL_Texture *)cmd->texture.ptr, (const float *)((const nk_byte *)vertices + vp), vs,
               (const SDL_FColor *)((const nk_byte *)vertices + vc), vs, (const float *)((const nk_byte *)vertices + vt), vs, (int)(vbuf.needed / vs),
               (void *)offset, (int)cmd->elem_count, 2);

            offset += cmd->elem_count;
         }
      }

      SDL_SetRenderClipRect(bridge->renderer, &saved_clip);
      if (!clipping_enabled)
      {
         SDL_SetRenderClipRect(bridge->renderer, NULL);
      }

      nk_clear(&bridge->ctx);
      nk_buffer_clear(&bridge->device.cmds);
      nk_buffer_free(&vbuf);
      nk_buffer_free(&ebuf);
   }
}

static void uiBridgeClipboardPaste(nk_handle const usr, struct nk_text_edit *const edit)
{
   NK_UNUSED(usr);

   char *const text = SDL_GetClipboardText();
   NK_ASSERT(text);

   if (text[0] != '\0')
   {
      int len;
      /* FIXME: there is a bug in Nuklear that affects UTF8 clipboard handling
       * "len" should be a buffer length, but due to bug it must be a glyph count
       * see: https://github.com/Immediate-Mode-UI/Nuklear/pull/841 */
#if 0
        len = nk_strlen(text);
#else
      len = (int)SDL_utf8strlen(text);
#endif
      nk_textedit_paste(edit, text, len);
   }
   SDL_free(text);
}

static void uiBridgeClipboardCopy(nk_handle const usr, char const *const text, int const len)
{
   if (len <= 0 || text == NULL)
      return;

   struct UIBridge const *bridge = usr.ptr;
   NK_ASSERT(bridge);

   /* FIXME: there is a bug in Nuklear that affects UTF8 clipboard handling
    * "len" is expected to be a buffer length, but due to bug it actually is a glyph count
    * see: https://github.com/Immediate-Mode-UI/Nuklear/pull/841 */
   size_t bufLen;
#if 0
   bufLen = len + 1;
#else
   char const *ptext = text;
   for (int i = len; i > 0; i--)
      (void)SDL_StepUTF8(&ptext, NULL);
   bufLen = (size_t)(ptext - text) + 1;
#endif

   char *const str = bridge->allocator.alloc(bridge->allocator.userdata, NULL, bufLen);
   if (!str)
      return;
   SDL_strlcpy(str, text, bufLen);
   SDL_SetClipboardText(str);
   bridge->allocator.free(bridge->allocator.userdata, str);
}

struct nk_context *uiBridgeInit(SDL_Window *const win, SDL_Renderer *const renderer)
{
   NK_ASSERT(win);
   NK_ASSERT(renderer);
   struct nk_allocator const allocator = nkAllocatorFromSDL();
   struct UIBridge *const bridge = allocator.alloc(allocator.userdata, NULL, sizeof(*bridge));
   NK_ASSERT(bridge);
   SDL_zerop(bridge);
   bridge->allocator = allocator;
   bridge->window = win;
   bridge->renderer = renderer;
   nk_init(&bridge->ctx, &bridge->allocator, NULL);
   bridge->ctx.userdata = nk_handle_ptr((void *)bridge);
   bridge->ctx.clip.copy = uiBridgeClipboardCopy;
   bridge->ctx.clip.paste = uiBridgeClipboardPaste;
   bridge->ctx.clip.userdata = nk_handle_ptr((void *)bridge);
   nk_buffer_init(&bridge->device.cmds, &bridge->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
   bridge->editWasActive = false;
   bridge->insertKeyToggle = false;
   bridge->lastRenderTick = SDL_GetTicksNS();
   return &bridge->ctx;
}

int uiBridgeHandleEvent(struct nk_context *const ctx, SDL_Event const *const evt)
{
   NK_ASSERT(ctx);
   NK_ASSERT(evt);

   struct UIBridge *const bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);

   /* We only care about Window currently used by Nuklear */
   if (bridge->window != SDL_GetWindowFromEvent(evt))
   {
      return 0;
   }

   switch (evt->type)
   {
   case SDL_EVENT_KEY_UP: /* KEYUP & KEYDOWN share same routine */
   case SDL_EVENT_KEY_DOWN:
   {
      bool down = evt->type == SDL_EVENT_KEY_DOWN;
      bool ctrl_down = (evt->key.mod & SDL_KMOD_CTRL) != 0;

      switch (evt->key.key)
      {
      case SDLK_LALT:
      case SDLK_RALT:
         nk_input_key(ctx, NK_KEY_ALT, down);
         break;
      case SDLK_RSHIFT: /* RSHIFT & LSHIFT share same routine */
      case SDLK_LSHIFT:
         nk_input_key(ctx, NK_KEY_SHIFT, down);
         break;
      case SDLK_DELETE:
         nk_input_key(ctx, NK_KEY_DEL, down);
         break;
      case SDLK_KP_ENTER:
      case SDLK_RETURN:
         nk_input_key(ctx, NK_KEY_ENTER, down);
         break;
      case SDLK_TAB:
         nk_input_key(ctx, NK_KEY_TAB, down);
         break;
      case SDLK_BACKSPACE:
         nk_input_key(ctx, NK_KEY_BACKSPACE, down);
         break;
      case SDLK_HOME:
         nk_input_key(ctx, NK_KEY_TEXT_START, down);
         nk_input_key(ctx, NK_KEY_SCROLL_START, down);
         break;
      case SDLK_END:
         nk_input_key(ctx, NK_KEY_TEXT_END, down);
         nk_input_key(ctx, NK_KEY_SCROLL_END, down);
         break;
      case SDLK_PAGEDOWN:
         nk_input_key(ctx, NK_KEY_SCROLL_DOWN, down);
         break;
      case SDLK_PAGEUP:
         nk_input_key(ctx, NK_KEY_SCROLL_UP, down);
         break;
      case SDLK_F1:
         nk_input_key(ctx, NK_KEY_F1, down);
         break;
      case SDLK_F2:
         nk_input_key(ctx, NK_KEY_F2, down);
         break;
      case SDLK_F3:
         nk_input_key(ctx, NK_KEY_F3, down);
         break;
      case SDLK_F4:
         nk_input_key(ctx, NK_KEY_F4, down);
         break;
      case SDLK_F5:
         nk_input_key(ctx, NK_KEY_F5, down);
         break;
      case SDLK_F6:
         nk_input_key(ctx, NK_KEY_F6, down);
         break;
      case SDLK_F7:
         nk_input_key(ctx, NK_KEY_F7, down);
         break;
      case SDLK_F8:
         nk_input_key(ctx, NK_KEY_F8, down);
         break;
      case SDLK_F9:
         nk_input_key(ctx, NK_KEY_F9, down);
         break;
      case SDLK_F10:
         nk_input_key(ctx, NK_KEY_F10, down);
         break;
      case SDLK_F11:
         nk_input_key(ctx, NK_KEY_F11, down);
         break;
      case SDLK_F12:
         nk_input_key(ctx, NK_KEY_F12, down);
         break;
      case SDLK_A:
         nk_input_key(ctx, NK_KEY_TEXT_SELECT_ALL, down && ctrl_down);
         break;
      case SDLK_Z:
         nk_input_key(ctx, NK_KEY_TEXT_UNDO, down && ctrl_down);
         break;
      case SDLK_R:
         nk_input_key(ctx, NK_KEY_TEXT_REDO, down && ctrl_down);
         break;
      case SDLK_C:
         nk_input_key(ctx, NK_KEY_COPY, down && ctrl_down);
         break;
      case SDLK_V:
         nk_input_key(ctx, NK_KEY_PASTE, down && ctrl_down);
         break;
      case SDLK_X:
         nk_input_key(ctx, NK_KEY_CUT, down && ctrl_down);
         break;
      case SDLK_B:
         nk_input_key(ctx, NK_KEY_TEXT_LINE_START, down && ctrl_down);
         break;
      case SDLK_E:
         nk_input_key(ctx, NK_KEY_TEXT_LINE_END, down && ctrl_down);
         break;
      case SDLK_UP:
         nk_input_key(ctx, NK_KEY_UP, down);
         break;
      case SDLK_DOWN:
         nk_input_key(ctx, NK_KEY_DOWN, down);
         break;
      case SDLK_ESCAPE:
         nk_input_key(ctx, NK_KEY_TEXT_RESET_MODE, down);
         break;
      case SDLK_INSERT:
         if (down)
            bridge->insertKeyToggle = !bridge->insertKeyToggle;
         if (bridge->insertKeyToggle)
         {
            nk_input_key(ctx, NK_KEY_TEXT_INSERT_MODE, down);
         }
         else
         {
            nk_input_key(ctx, NK_KEY_TEXT_REPLACE_MODE, down);
         }
         break;
      case SDLK_LEFT:
         if (ctrl_down)
            nk_input_key(ctx, NK_KEY_TEXT_WORD_LEFT, down);
         else
            nk_input_key(ctx, NK_KEY_LEFT, down);
         break;
      case SDLK_RIGHT:
         if (ctrl_down)
            nk_input_key(ctx, NK_KEY_TEXT_WORD_RIGHT, down);
         else
            nk_input_key(ctx, NK_KEY_RIGHT, down);
         break;
      default:
         return 0;
      }
      return 1;
   }

   case SDL_EVENT_MOUSE_BUTTON_UP: /* MOUSEBUTTONUP & MOUSEBUTTONDOWN share same routine */
   case SDL_EVENT_MOUSE_BUTTON_DOWN:
   {
      int const x = (int)nk_roundf(evt->button.x);
      int const y = (int)nk_roundf(evt->button.y);
      bool const down = evt->button.down != 0;
      switch (evt->button.button)
      {
      case SDL_BUTTON_LEFT:
         if (evt->button.clicks > 1)
            nk_input_button(ctx, NK_BUTTON_DOUBLE, x, y, down);
         nk_input_button(ctx, NK_BUTTON_LEFT, x, y, down);
         break;
      case SDL_BUTTON_MIDDLE:
         nk_input_button(ctx, NK_BUTTON_MIDDLE, x, y, down);
         break;
      case SDL_BUTTON_RIGHT:
         nk_input_button(ctx, NK_BUTTON_RIGHT, x, y, down);
         break;
      case SDL_BUTTON_X1:
         nk_input_button(ctx, NK_BUTTON_X1, x, y, down);
         break;
      case SDL_BUTTON_X2:
         nk_input_button(ctx, NK_BUTTON_X2, x, y, down);
         break;
      default:
         return 0;
      }
   }
      return 1;

   case SDL_EVENT_MOUSE_MOTION:
      ctx->input.mouse.pos.x = evt->motion.x;
      ctx->input.mouse.pos.y = evt->motion.y;
      ctx->input.mouse.delta.x = ctx->input.mouse.pos.x - ctx->input.mouse.prev.x;
      ctx->input.mouse.delta.y = ctx->input.mouse.pos.y - ctx->input.mouse.prev.y;
      return 1;

   case SDL_EVENT_TEXT_INPUT:
   {
      NK_ASSERT(evt->text.text);
      nk_size len = SDL_strlen(evt->text.text);
      NK_ASSERT(len <= NK_UTF_SIZE);
      nk_glyph glyph;
      NK_MEMCPY(glyph, evt->text.text, len);
      nk_input_glyph(ctx, glyph);
   }
      return 1;

   case SDL_EVENT_MOUSE_WHEEL:
      nk_input_scroll(ctx, nk_vec2(evt->wheel.x, evt->wheel.y));
      return 1;

   default:
      return 0;
   }
   return 0;
}

void uiBridgeShutdown(struct nk_context *const ctx)
{
   NK_ASSERT(ctx);
   struct UIBridge *const bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);

   nk_buffer_free(&bridge->device.cmds);

   uiBridgeDeviceDropFontTexture(&bridge->device);
   fontRelease(&bridge->uiFont);

   nk_free(ctx);
   bridge->allocator.free(bridge->allocator.userdata, bridge->nkFont);
   bridge->allocator.free(bridge->allocator.userdata, bridge);
}

static float uiBridgeQueryFontWidth(nk_handle const handle, float const height, char const *text, int const len)
{
   struct UIBridge const *const bridge = handle.ptr;
   NK_ASSERT(bridge);
   int32_t width = 0;

   nk_rune unicode = 0;
   int offset = 0;
   int charsUsed = nk_utf_decode(text, &unicode, len);
   while ((charsUsed > 0) && (offset < len))
   {
      FontCodepointEntry const *entry = fontFindCodepointEntry(bridge->uiFont, unicode);
      if (entry == NULL)
      {
         entry = fontFindCodepointEntry(bridge->uiFont, '?');
      }
      width += entry->rect.size.width - 1; // unclear why it is -1 compared to how the query function operates; yet it works.

      offset += charsUsed;
      charsUsed = nk_utf_decode(&text[offset], &unicode, len - offset);
   }
   float const scale = height / (float)(bridge->uiFont->height + 2);
   return (float)(width + 1) * scale;
}

static void uiBridgeQueryFontGlyph(
   nk_handle const handle, float const height, struct nk_user_font_glyph *const glyph, nk_rune const codepoint, nk_rune const next_codepoint)
{
   NK_UNUSED(next_codepoint);

   struct UIBridge const *const bridge = handle.ptr;
   NK_ASSERT(bridge);

   FontCodepointEntry const *entry = fontFindCodepointEntry(bridge->uiFont, codepoint);
   if (entry == NULL)
   {
      entry = fontFindCodepointEntry(bridge->uiFont, '?');
   }

   float const bitmapWidth = bridge->uiFont->atlas.size.width;
   float const bitmapHeight = bridge->uiFont->atlas.size.height;
   int const outlinedLeft = entry->rect.topLeft.x - 1;
   int const outlinedTop = entry->rect.topLeft.y - 1;
   int const outlinedHeight = entry->rect.size.height + 2;
   int const outlinedWidth = entry->rect.size.width + 2;
   float const scale = height / (float)(outlinedHeight);
   glyph->height = (float)(outlinedHeight)*scale;
   glyph->width = (float)(outlinedWidth)*scale;
   glyph->xadvance = (float)(outlinedWidth - 2) * scale; // one shift for outline, one because each glyph has its own spacing to the right
   // The renderer on MS-DOS appears to have a different rounding or some rounding error. This can be fixed by adding float's epsilon with no effect on others.
   static float const dosFudge = SDL_FLT_EPSILON;
   glyph->uv[0].x = (float)(outlinedLeft) / bitmapWidth + dosFudge;
   glyph->uv[0].y = (float)(outlinedTop) / bitmapHeight + dosFudge;
   glyph->uv[1].x = (float)(outlinedLeft + outlinedWidth) / bitmapWidth + (dosFudge * 2.0f);
   glyph->uv[1].y = (float)(outlinedTop + outlinedHeight) / bitmapHeight + (dosFudge * 2.0f);
   glyph->offset.x = 0.0f;
   glyph->offset.y = 0.0f;
}

static void renderMonochromeFontCharacter(
   SDL_Surface *const surface, size_t const characterIndex, PixelOffset const off, Font const *const font, struct nk_color const color)
{
   FontCodepointEntry const *const entry = &font->codepoints[characterIndex];
   PixelAxisOffset const outXBase = (PixelAxisOffset)entry->rect.topLeft.x + off.x;
   PixelAxisOffset const outYBase = (PixelAxisOffset)entry->rect.topLeft.y + off.y;
   for (PixelAxisSize y = 0; y < entry->rect.size.height; ++y)
   {
      for (PixelAxisSize x = 0; x < entry->rect.size.width; ++x)
      {
         uint8_t const in = *(font->atlas.data + (font->atlas.stride * (entry->rect.topLeft.y + y)) + entry->rect.topLeft.x + x);
         if (in != 0)
         {
            SDL_WriteSurfacePixel(surface, outXBase + x, outYBase + y, color.r, color.g, color.b, color.a);
         }
      }
   }
}

void uiBridgeSetFont(struct nk_context *const ctx, float const scale)
{
   NK_ASSERT(ctx);

   struct UIBridge *const bridge = ctx->userdata.ptr;
   NK_ASSERT(bridge);

   if (bridge->nkFont)
   {
      bridge->allocator.free(bridge->allocator.userdata, bridge->nkFont);
      bridge->nkFont = NULL;
   }

   bridge->uiFont = uiFont();

   SDL_Surface *const surface = SDL_CreateSurface(bridge->uiFont->atlas.size.width, bridge->uiFont->atlas.size.height, SDL_PIXELFORMAT_RGBA32);
   NK_ASSERT(surface);

   struct nk_color const black = {.r = 0x30, .g = 0x30, .b = 0x30, .a = 0xFF};

   for (size_t currentCharOffset = 0; currentCharOffset < bridge->uiFont->codepointCount; ++currentCharOffset)
   {
      for (PixelAxisOffset i = 0; i < 9; ++i)
      {
         PixelOffset const off = {.x = (PixelAxisOffset)((i % 3) - 1), .y = (PixelAxisOffset)((i / 3) - 1)};
         if (off.x != 0 || off.y != 0)
         {
            renderMonochromeFontCharacter(surface, currentCharOffset, off, bridge->uiFont, black);
         }
      }
      struct nk_color const textColor = {.r = 0xFF, .g = 0xFF, .b = 0xFF, .a = 0xFF};
      PixelOffset const zeroOffset = {};
      renderMonochromeFontCharacter(surface, currentCharOffset, zeroOffset, bridge->uiFont, textColor);
   }

   struct nk_user_font *const font = bridge->allocator.alloc(bridge->allocator.userdata, NULL, sizeof(*font));
   NK_ASSERT(font);
   font->userdata.ptr = bridge;
   font->height = (float)(bridge->uiFont->height + 2) * scale;
   font->width = &uiBridgeQueryFontWidth;
   font->query = &uiBridgeQueryFontGlyph;

   /* HACK: upload atlas turns pixels into SDL_Texture
    *       and sets said Texture into sdl->ogl.font_tex
    *       then nk_sdl_render expects same Texture at font->texture */
   uiBridgeUploadAtlas(ctx, surface->pixels, surface->w, surface->h);
   font->texture.ptr = bridge->device.fontTexture;

   bridge->nkFont = font;
   nk_style_set_font(ctx, font);

   SDL_DestroySurface(surface);
}
