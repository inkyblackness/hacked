#include <stddef.h>

#include "hacked/core/editor/ui/UIFont.h"
#include "hacked/core/media/PixelSpace.h"
#include "hacked/core/serial/io/Primitives.h"

#define NK_IMPLEMENTATION
#include "hacked/nuklear/NuklearSdlBridge.h"

// This is the same default value as the one from "src/nuklear_internal.h"
#ifndef NK_BUFFER_DEFAULT_INITIAL_SIZE
#define NK_BUFFER_DEFAULT_INITIAL_SIZE (4 * 1024)
#endif

struct nk_sdl_device
{
   struct nk_buffer cmds;
   struct nk_draw_null_texture tex_null;
   SDL_Texture *font_tex;
};

struct nk_sdl_vertex
{
   float position[2];
   float uv[2];
   float col[4];
};

struct nk_sdl
{
   SDL_Window *win;
   SDL_Renderer *renderer;
   Font const *uiFont;
   struct nk_user_font *debug_font;
   struct nk_sdl_device ogl;
   struct nk_context ctx;
   struct nk_allocator allocator;
   nk_handle userdata;
   Uint64 last_render;
   bool insert_toggle;
   bool edit_was_active;
};

NK_API nk_handle nk_sdl_get_userdata(struct nk_context *ctx)
{
   NK_ASSERT(ctx);
   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);
   return sdl->userdata;
}

NK_API void nk_sdl_set_userdata(struct nk_context *ctx, nk_handle userdata)
{
   NK_ASSERT(ctx);
   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);
   sdl->userdata = userdata;
}

NK_INTERN void *nk_sdl_alloc(nk_handle user, void *old, nk_size size)
{
   NK_UNUSED(user);
   /* FIXME: nk_sdl_alloc should use SDL_realloc here, not SDL_malloc
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

NK_INTERN void nk_sdl_free(nk_handle user, void *old)
{
   NK_UNUSED(user);
   SDL_free(old);
}

NK_API struct nk_allocator nk_sdl_allocator()
{
   struct nk_allocator allocator;
   allocator.userdata.ptr = NULL;
   allocator.alloc = nk_sdl_alloc;
   allocator.free = nk_sdl_free;
   return allocator;
}

NK_INTERN void nk_sdl_device_upload_atlas(struct nk_context *ctx, void const *image, int width, int height)
{
   NK_ASSERT(ctx);
   NK_ASSERT(image);
   NK_ASSERT(width > 0);
   NK_ASSERT(height > 0);

   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   /* Clean up if the texture already exists. */
   if (sdl->ogl.font_tex != NULL)
   {
      SDL_DestroyTexture(sdl->ogl.font_tex);
      sdl->ogl.font_tex = NULL;
   }

   sdl->ogl.font_tex = SDL_CreateTexture(sdl->renderer, SDL_PIXELFORMAT_ARGB8888, SDL_TEXTUREACCESS_STATIC, width, height);
   NK_ASSERT(sdl->ogl.font_tex);
   SDL_UpdateTexture(sdl->ogl.font_tex, NULL, image, 4 * width);
   SDL_SetTextureBlendMode(sdl->ogl.font_tex, SDL_BLENDMODE_BLEND);
   SDL_SetTextureScaleMode(sdl->ogl.font_tex, SDL_SCALEMODE_NEAREST);
}

NK_API void nk_sdl_update_TextInput(struct nk_context *ctx)
{
   bool active = false;
   NK_ASSERT(ctx);
   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

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
   if (active != sdl->edit_was_active)
   {
      bool const window_edit_active = SDL_TextInputActive(sdl->win);

      /* If you ever hit this check, it means that the demo and your app
       * (or something else) are all trying to manage TextInputActive state.
       * This can cause subtle bugs where the state won't be what you expect.
       * You can safely remove this assert and the demo will keep working,
       * but make sure it does not cause any issues for you */
      NK_ASSERT(window_edit_active == sdl->edit_was_active && "something else changed TextInputActive state for this Window");

      if (!window_edit_active && !sdl->edit_was_active && active)
         SDL_StartTextInput(sdl->win);
      else if (window_edit_active && sdl->edit_was_active && !active)
         SDL_StopTextInput(sdl->win);
      sdl->edit_was_active = active;
   }

   /* FIXME:
    * for full SDL3 integration, you also need to find current edit widget
    * bounds and the text cursor offset, and pass this data into SDL_SetTextInputArea.
    * This is currently not possible to do safely as Nuklear does not support it.
    * https://wiki.libsdl.org/SDL3/SDL_SetTextInputArea
    * https://github.com/Immediate-Mode-UI/Nuklear/pull/857
    */
}

NK_API void nk_sdl_render(struct nk_context *ctx, enum nk_anti_aliasing AA)
{
   NK_ASSERT(ctx);
   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   { /* setup internal delta time that Nuklear needs for animations */
      Uint64 const now = SDL_GetTicksNS();
      ctx->delta_time_seconds = (float)(now - sdl->last_render) / (float)SDL_NS_PER_SECOND;
      sdl->last_render = now;
   }

   {
      int const vs = sizeof(struct nk_sdl_vertex);
      size_t const vp = NK_OFFSETOF(struct nk_sdl_vertex, position);
      size_t const vt = NK_OFFSETOF(struct nk_sdl_vertex, uv);
      size_t const vc = NK_OFFSETOF(struct nk_sdl_vertex, col);

      /* convert from command queue into draw list and draw to screen */

      /* fill converting configuration */
      static const struct nk_draw_vertex_layout_element vertex_layout[] = {{NK_VERTEX_POSITION, NK_FORMAT_FLOAT, NK_OFFSETOF(struct nk_sdl_vertex, position)},
         {NK_VERTEX_TEXCOORD, NK_FORMAT_FLOAT, NK_OFFSETOF(struct nk_sdl_vertex, uv)},
         {NK_VERTEX_COLOR, NK_FORMAT_R32G32B32A32_FLOAT, NK_OFFSETOF(struct nk_sdl_vertex, col)}, {NK_VERTEX_LAYOUT_END}};
      struct nk_convert_config config;
      NK_MEMSET(&config, 0, sizeof(config));
      config.vertex_layout = vertex_layout;
      config.vertex_size = sizeof(struct nk_sdl_vertex);
      config.vertex_alignment = NK_ALIGNOF(struct nk_sdl_vertex);
      config.tex_null = sdl->ogl.tex_null;
      config.circle_segment_count = 22;
      config.curve_segment_count = 22;
      config.arc_segment_count = 22;
      config.global_alpha = 1.0f;
      config.shape_AA = AA;
      config.line_AA = AA;

      /* convert shapes into vertices */
      struct nk_buffer vbuf, ebuf;
      nk_buffer_init(&vbuf, &sdl->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
      nk_buffer_init(&ebuf, &sdl->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
      nk_convert(&sdl->ctx, &sdl->ogl.cmds, &vbuf, &ebuf, &config);

      /* iterate over and execute each draw command */
      nk_draw_index const *offset = (const nk_draw_index *)nk_buffer_memory_const(&ebuf);

      bool clipping_enabled = SDL_RenderClipEnabled(sdl->renderer);
      SDL_Rect saved_clip;
      SDL_GetRenderClipRect(sdl->renderer, &saved_clip);

      struct nk_draw_command const *cmd = NULL;
      nk_draw_foreach(cmd, &sdl->ctx, &sdl->ogl.cmds)
      {
         if (!cmd->elem_count)
            continue;

         {
            SDL_Rect r;
            r.x = (int)nk_roundf(cmd->clip_rect.x);
            r.y = (int)nk_roundf(cmd->clip_rect.y);
            r.w = (int)nk_roundf(cmd->clip_rect.w);
            r.h = (int)nk_roundf(cmd->clip_rect.h);
            SDL_SetRenderClipRect(sdl->renderer, &r);
         }

         {
            const void *vertices = nk_buffer_memory_const(&vbuf);

            SDL_RenderGeometryRaw(sdl->renderer, (SDL_Texture *)cmd->texture.ptr, (const float *)((const nk_byte *)vertices + vp), vs,
               (const SDL_FColor *)((const nk_byte *)vertices + vc), vs, (const float *)((const nk_byte *)vertices + vt), vs, (int)(vbuf.needed / vs),
               (void *)offset, (int)cmd->elem_count, 2);

            offset += cmd->elem_count;
         }
      }

      SDL_SetRenderClipRect(sdl->renderer, &saved_clip);
      if (!clipping_enabled)
      {
         SDL_SetRenderClipRect(sdl->renderer, NULL);
      }

      nk_clear(&sdl->ctx);
      nk_buffer_clear(&sdl->ogl.cmds);
      nk_buffer_free(&vbuf);
      nk_buffer_free(&ebuf);
   }
}

NK_INTERN void nk_sdl_clipboard_paste(nk_handle usr, struct nk_text_edit *edit)
{
   NK_UNUSED(usr);

   /* this function returns empty string on failure, not NULL */
   char *text = SDL_GetClipboardText();
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

NK_INTERN void nk_sdl_clipboard_copy(nk_handle usr, char const *text, int len)
{
   if (len <= 0 || text == NULL)
      return;

   struct nk_sdl const *sdl = (struct nk_sdl *)usr.ptr;
   NK_ASSERT(sdl);

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

   char *str = (char *)sdl->allocator.alloc(sdl->allocator.userdata, NULL, bufLen);
   if (!str)
      return;
   SDL_strlcpy(str, text, bufLen);
   SDL_SetClipboardText(str);
   sdl->allocator.free(sdl->allocator.userdata, str);
}

NK_API struct nk_context *nk_sdl_init(SDL_Window *win, SDL_Renderer *renderer, struct nk_allocator allocator)
{
   NK_ASSERT(win);
   NK_ASSERT(renderer);
   NK_ASSERT(allocator.alloc);
   NK_ASSERT(allocator.free);
   struct nk_sdl *sdl = (struct nk_sdl *)allocator.alloc(allocator.userdata, NULL, sizeof(*sdl));
   NK_ASSERT(sdl);
   SDL_zerop(sdl);
   sdl->allocator.userdata = allocator.userdata;
   sdl->allocator.alloc = allocator.alloc;
   sdl->allocator.free = allocator.free;
   sdl->win = win;
   sdl->renderer = renderer;
   nk_init(&sdl->ctx, &sdl->allocator, NULL);
   sdl->ctx.userdata = nk_handle_ptr((void *)sdl);
   sdl->ctx.clip.copy = nk_sdl_clipboard_copy;
   sdl->ctx.clip.paste = nk_sdl_clipboard_paste;
   sdl->ctx.clip.userdata = nk_handle_ptr((void *)sdl);
   nk_buffer_init(&sdl->ogl.cmds, &sdl->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
   sdl->edit_was_active = false;
   sdl->insert_toggle = false;
   sdl->last_render = SDL_GetTicksNS();
   return &sdl->ctx;
}

NK_API int nk_sdl_handle_event(struct nk_context *ctx, SDL_Event *evt)
{
   NK_ASSERT(ctx);
   NK_ASSERT(evt);

   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   /* We only care about Window currently used by Nuklear */
   if (sdl->win != SDL_GetWindowFromEvent(evt))
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
            sdl->insert_toggle = !sdl->insert_toggle;
         if (sdl->insert_toggle)
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

NK_API
void nk_sdl_shutdown(struct nk_context *ctx)
{
   NK_ASSERT(ctx);
   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   nk_buffer_free(&sdl->ogl.cmds);

   if (sdl->ogl.font_tex != NULL)
   {
      SDL_DestroyTexture(sdl->ogl.font_tex);
      sdl->ogl.font_tex = NULL;
   }

   nk_free(ctx);
   sdl->allocator.free(sdl->allocator.userdata, sdl->debug_font);
   sdl->allocator.free(sdl->allocator.userdata, sdl);
   fontRelease(&sdl->uiFont);
}

NK_INTERN float nk_sdl_query_tiny_font_width(nk_handle const handle, float const height, char const *text, int const len)
{
   struct nk_sdl *sdl = (struct nk_sdl *)handle.ptr;
   NK_ASSERT(sdl);
   int32_t width = 0;
   char const *const end = text + len;
   for (char const *it = text; it != end; ++it)
   {
      // TODO this is wrong, as the text is UTF-8, and it needs to be decoded and then translated into the right codepage.
      char codepoint = *it;

      FontCodepointEntry const *entry = fontFindCodepointEntry(sdl->uiFont, codepoint);
      if (entry == NULL)
      {
         entry = fontFindCodepointEntry(sdl->uiFont, '?');
      }

      width += entry->rect.size.width - 1; // unclear why it is -1 compared to how the query function operates; yet it works.
   }
   // TODO: consider having a reference height value in the font structure
   float scale = height / (float)(sdl->uiFont->codepoints[0].rect.size.height + 2);
   return (float)(width + 1) * scale;
}

NK_INTERN void nk_sdl_query_tiny_font_glyph(
   nk_handle const handle, float const height, struct nk_user_font_glyph *const glyph, nk_rune const codepoint, nk_rune const next_codepoint)
{
   NK_UNUSED(next_codepoint);

   struct nk_sdl *sdl = (struct nk_sdl *)handle.ptr;
   NK_ASSERT(sdl);

   FontCodepointEntry const *entry = fontFindCodepointEntry(sdl->uiFont, codepoint);
   if (entry == NULL)
   {
      entry = fontFindCodepointEntry(sdl->uiFont, '?');
   }

   float const bitmapWidth = sdl->uiFont->atlas.size.width;
   float const bitmapHeight = sdl->uiFont->atlas.size.height;
   int const outlinedLeft = entry->rect.topLeft.x - 1;
   int const outlinedTop = entry->rect.topLeft.y - 1;
   int const outlinedHeight = entry->rect.size.height + 2;
   int const outlinedWidth = entry->rect.size.width + 2;
   float const scale = height / (float)(outlinedHeight);
   glyph->height = (float)(outlinedHeight)*scale;
   glyph->width = (float)(outlinedWidth)*scale;
   glyph->xadvance = (float)(outlinedWidth - 2) * scale; // one shift for outline, one because each glyph has its own spacing to the right
   glyph->uv[0].x = (float)(outlinedLeft) / bitmapWidth;
   glyph->uv[0].y = (float)(outlinedTop) / bitmapHeight;
   glyph->uv[1].x = (float)(outlinedLeft + outlinedWidth) / bitmapWidth;
   glyph->uv[1].y = (float)(outlinedTop + outlinedHeight) / bitmapHeight;
   glyph->offset.x = 0.0f;
   glyph->offset.y = 0.0f;
}

typedef struct
{
   uint8_t r;
   uint8_t g;
   uint8_t b;
} ColorRgb;

static void renderMonochromeFontCharacter(
   SDL_Surface *const surface, size_t const characterIndex, PixelOffset const off, Font const *const font, ColorRgb const color)
{
   // TODO figure out types; when which is more appropriate.
   FontCodepointEntry const *const entry = &font->codepoints[characterIndex];
   int outXBase = entry->rect.topLeft.x + off.x;
   int outYBase = entry->rect.topLeft.y + off.y;
   for (int16_t y = 0; y < entry->rect.size.height; ++y)
   {
      for (int16_t x = 0; x < entry->rect.size.width; ++x)
      {
         uint8_t const in = *(font->atlas.data + (font->atlas.stride * (entry->rect.topLeft.y + y)) + entry->rect.topLeft.x + x);
         if (in != 0)
         {
            SDL_WriteSurfacePixel(surface, outXBase + x, outYBase + y, color.b, color.g, color.r, 0xFF); // TODO: Why is this BGRA and not RGBA?
         }
      }
   }
}

NK_API void nk_sdl_style_set_tiny_font(struct nk_context *const ctx, float const scale)
{
   NK_ASSERT(ctx);

   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   if (sdl->debug_font)
   {
      sdl->allocator.free(sdl->allocator.userdata, sdl->debug_font);
      sdl->debug_font = NULL;
   }

   sdl->uiFont = uiFont();

   SDL_Surface *surface = SDL_CreateSurface(sdl->uiFont->atlas.size.width, sdl->uiFont->atlas.size.height, SDL_PIXELFORMAT_RGBA32);
   NK_ASSERT(surface);

   ColorRgb black = {};

   for (size_t currentCharOffset = 0; currentCharOffset < sdl->uiFont->codepointCount; ++currentCharOffset)
   {
      for (PixelAxisOffset i = 0; i < 9; ++i)
      {
         PixelOffset const off = {.x = (PixelAxisOffset)((i % 3) - 1), .y = (PixelAxisOffset)((i / 3) - 1)};
         if (off.x != 0 || off.y != 0)
         {
            renderMonochromeFontCharacter(surface, currentCharOffset, off, sdl->uiFont, black);
         }
      }
      ColorRgb textColor = {.r = 0x5B, .g = 0xAC, .b = 0x1E};
      PixelOffset zeroOffset = {};
      renderMonochromeFontCharacter(surface, currentCharOffset, zeroOffset, sdl->uiFont, textColor);
   }

   struct nk_user_font *font = (struct nk_user_font *)sdl->allocator.alloc(sdl->allocator.userdata, NULL, sizeof(*font));
   NK_ASSERT(font);
   font->userdata.ptr = sdl;
   font->height = (float)(sdl->uiFont->codepoints[0].rect.size.height + 2) * scale;
   font->width = &nk_sdl_query_tiny_font_width;
   font->query = &nk_sdl_query_tiny_font_glyph;

   /* HACK: nk_sdl_device_upload_atlas turns pixels into SDL_Texture
    *       and sets said Texture into sdl->ogl.font_tex
    *       then nk_sdl_render expects same Texture at font->texture */
   nk_sdl_device_upload_atlas(ctx, surface->pixels, surface->w, surface->h);
   font->texture.ptr = sdl->ogl.font_tex;

   sdl->debug_font = font;
   nk_style_set_font(ctx, font);

   SDL_DestroySurface(surface);
}
