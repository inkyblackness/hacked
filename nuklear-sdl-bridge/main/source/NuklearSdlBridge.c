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
   struct nk_user_font *debug_font;
   struct nk_sdl_device ogl;
   struct nk_context ctx;
#ifdef NK_INCLUDE_FONT_BAKING
   struct nk_font_atlas atlas;
#endif
   struct nk_allocator allocator;
   nk_handle userdata;
   Uint64 last_render;
   bool insert_toggle;
   bool edit_was_active;
};

static uint8_t const tinyFont[] = {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x54,
   0x00, 0x00, 0x00, 0x16, 0x02, 0x00, 0x00, 0x4D, 0x00, 0x06, 0x00, 0x00, 0x00, 0x02, 0x00, 0x04, 0x00, 0x08, 0x00, 0x0C, 0x00, 0x10, 0x00, 0x14, 0x00, 0x18,
   0x00, 0x1A, 0x00, 0x1D, 0x00, 0x20, 0x00, 0x24, 0x00, 0x28, 0x00, 0x2B, 0x00, 0x2F, 0x00, 0x31, 0x00, 0x35, 0x00, 0x39, 0x00, 0x3D, 0x00, 0x41, 0x00, 0x45,
   0x00, 0x49, 0x00, 0x4D, 0x00, 0x51, 0x00, 0x55, 0x00, 0x59, 0x00, 0x5D, 0x00, 0x5F, 0x00, 0x62, 0x00, 0x66, 0x00, 0x69, 0x00, 0x6D, 0x00, 0x71, 0x00, 0x75,
   0x00, 0x79, 0x00, 0x7D, 0x00, 0x81, 0x00, 0x85, 0x00, 0x89, 0x00, 0x8D, 0x00, 0x91, 0x00, 0x95, 0x00, 0x97, 0x00, 0x9B, 0x00, 0x9F, 0x00, 0xA3, 0x00, 0xA9,
   0x00, 0xAE, 0x00, 0xB2, 0x00, 0xB6, 0x00, 0xBC, 0x00, 0xC0, 0x00, 0xC4, 0x00, 0xC8, 0x00, 0xCC, 0x00, 0xD0, 0x00, 0xD6, 0x00, 0xDA, 0x00, 0xDE, 0x00, 0xE2,
   0x00, 0xE5, 0x00, 0xE9, 0x00, 0xEC, 0x00, 0xF0, 0x00, 0xF4, 0x00, 0xF7, 0x00, 0xFB, 0x00, 0xFF, 0x00, 0x03, 0x01, 0x07, 0x01, 0x0B, 0x01, 0x0F, 0x01, 0x13,
   0x01, 0x17, 0x01, 0x19, 0x01, 0x1D, 0x01, 0x21, 0x01, 0x23, 0x01, 0x29, 0x01, 0x2D, 0x01, 0x31, 0x01, 0x35, 0x01, 0x3A, 0x01, 0x3E, 0x01, 0x42, 0x01, 0x46,
   0x01, 0x4A, 0x01, 0x4E, 0x01, 0x54, 0x01, 0x58, 0x01, 0x5C, 0x01, 0x60, 0x01, 0x64, 0x01, 0x68, 0x01, 0x6C, 0x01, 0x70, 0x01, 0x74, 0x01, 0x78, 0x01, 0x7C,
   0x01, 0x80, 0x01, 0x84, 0x01, 0x88, 0x01, 0x8C, 0x01, 0x90, 0x01, 0x94, 0x01, 0x98, 0x01, 0x9C, 0x01, 0xA0, 0x01, 0xA4, 0x01, 0xA8, 0x01, 0xAB, 0x01, 0xAF,
   0x01, 0xB3, 0x01, 0xB7, 0x01, 0xBB, 0x01, 0xBF, 0x01, 0xC3, 0x01, 0xC7, 0x01, 0xCB, 0x01, 0xCF, 0x01, 0xD3, 0x01, 0xD7, 0x01, 0xDB, 0x01, 0xDF, 0x01, 0xE3,
   0x01, 0xE7, 0x01, 0xEB, 0x01, 0xEF, 0x01, 0xF3, 0x01, 0xF7, 0x01, 0xFA, 0x01, 0xFE, 0x01, 0x02, 0x02, 0x06, 0x02, 0x0B, 0x02, 0x0C, 0x02, 0x0D, 0x02, 0x0E,
   0x02, 0x0F, 0x02, 0x10, 0x02, 0x11, 0x02, 0x12, 0x02, 0x13, 0x02, 0x14, 0x02, 0x15, 0x02, 0x16, 0x02, 0x17, 0x02, 0x18, 0x02, 0x19, 0x02, 0x1A, 0x02, 0x1B,
   0x02, 0x1C, 0x02, 0x1D, 0x02, 0x1E, 0x02, 0x1F, 0x02, 0x20, 0x02, 0x21, 0x02, 0x22, 0x02, 0x23, 0x02, 0x24, 0x02, 0x25, 0x02, 0x26, 0x02, 0x27, 0x02, 0x28,
   0x02, 0x29, 0x02, 0x2A, 0x02, 0x2B, 0x02, 0x2C, 0x02, 0x2D, 0x02, 0x2E, 0x02, 0x2F, 0x02, 0x30, 0x02, 0x31, 0x02, 0x32, 0x02, 0x33, 0x02, 0x34, 0x02, 0x35,
   0x02, 0x36, 0x02, 0x37, 0x02, 0x38, 0x02, 0x39, 0x02, 0x3A, 0x02, 0x3B, 0x02, 0x3C, 0x02, 0x3D, 0x02, 0x3E, 0x02, 0x3F, 0x02, 0x40, 0x02, 0x41, 0x02, 0x42,
   0x02, 0x43, 0x02, 0x44, 0x02, 0x45, 0x02, 0x46, 0x02, 0x4A, 0x02, 0x4B, 0x02, 0x4C, 0x02, 0x4D, 0x02, 0x4E, 0x02, 0x4F, 0x02, 0x50, 0x02, 0x51, 0x02, 0x52,
   0x02, 0x53, 0x02, 0x54, 0x02, 0x55, 0x02, 0x56, 0x02, 0x57, 0x02, 0x58, 0x02, 0x59, 0x02, 0x5A, 0x02, 0x5B, 0x02, 0x5C, 0x02, 0x5D, 0x02, 0x5E, 0x02, 0x5F,
   0x02, 0x60, 0x02, 0x61, 0x02, 0x62, 0x02, 0x63, 0x02, 0x64, 0x02, 0x65, 0x02, 0x66, 0x02, 0x67, 0x02, 0x68, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA3, 0x2A, 0xC4, 0x02, 0xAC, 0x04, 0x94, 0x86, 0x00, 0x95, 0x89, 0x95,
   0x54, 0x00, 0x00, 0x06, 0x8C, 0xD9, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2A, 0x4E, 0xD6, 0x92, 0xA4, 0x00, 0x17, 0x67,
   0x75, 0x74, 0x77, 0x74, 0x80, 0x03, 0x37, 0x67, 0x67, 0x77, 0x55, 0xD5, 0x11, 0x4B, 0xB3, 0xCC, 0xEE, 0xAA, 0x8A, 0xAB, 0xB4, 0x64, 0x08, 0x10, 0x04, 0x0C,
   0x11, 0x24, 0x40, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0E, 0x00, 0x50, 0x00, 0x05, 0x00, 0xAA, 0x49, 0xDC, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x80, 0x00, 0x01, 0x00, 0x02, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x2A, 0xEC, 0xE6, 0xA1, 0x4E, 0x1C, 0x25, 0x21, 0x17,
   0x47, 0x12, 0x50, 0x13, 0x21, 0x55, 0x74, 0x56, 0x64, 0x54, 0x99, 0x1B, 0x6A, 0xAA, 0x4A, 0x84, 0xAA, 0xA9, 0x28, 0xA2, 0x2A, 0x04, 0xD9, 0xCD, 0x89, 0xDC,
   0x05, 0x5F, 0x77, 0x77, 0x39, 0xBA, 0xAA, 0xAA, 0xAC, 0x00, 0x00, 0x08, 0xAC, 0x66, 0x66, 0xEC, 0xCC, 0x00, 0x15, 0x58, 0x01, 0xDD, 0xD5, 0x55, 0xD5, 0xC0,
   0x00, 0x0C, 0x3A, 0xBB, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x80, 0x00, 0x00, 0x00, 0x00, 0x46, 0x74, 0x21, 0xA4, 0x00, 0x25, 0x26, 0x31, 0x35,
   0x15, 0x74, 0xA0, 0x12, 0x37, 0x54, 0x54, 0x45, 0x74, 0x95, 0x15, 0x5A, 0xB2, 0xCC, 0x64, 0xAA, 0xA9, 0x11, 0x22, 0x20, 0x01, 0x55, 0x15, 0x9D, 0x55, 0x26,
   0x55, 0x55, 0x55, 0x21, 0x12, 0xAA, 0xA4, 0xA4, 0x00, 0x00, 0x08, 0xAC, 0xAA, 0xAA, 0x8C, 0xCC, 0x44, 0x5D, 0xD0, 0x01, 0x55, 0x55, 0x55, 0x55, 0x00, 0x00,
   0x15, 0x2A, 0xAA, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x80, 0x00, 0x00, 0x00, 0x20, 0xEE, 0xBA, 0x12, 0x00, 0x41, 0x47, 0x77, 0x71, 0x77, 0x17,
   0x11, 0x13, 0x20, 0x35, 0x77, 0x67, 0x47, 0x55, 0x95, 0xD1, 0x4B, 0xA3, 0xEA, 0xE4, 0xE4, 0x52, 0x93, 0xB1, 0x60, 0xE1, 0xDD, 0xDD, 0xC9, 0xD5, 0x65, 0x55,
   0x57, 0x77, 0x23, 0x13, 0x91, 0x4A, 0xE6, 0x00, 0x00, 0x0E, 0xEE, 0xEE, 0xEE, 0xEE, 0xEE, 0x44, 0x55, 0x5C, 0x01, 0xDD, 0xDD, 0xDD, 0xDD, 0xC0, 0x00, 0x1D,
   0x3B, 0xAA, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x80, 0x00, 0x00, 0x00, 0x00, 0x44, 0x04, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x00, 0x00, 0x00,
   0x41, 0x80, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x80, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00};

static int16_t tinyFontCharLowest = 0;
static int16_t tinyFontCharHighest = 0;
static int16_t tinyFontHeight = 0;
static int16_t *tinyFontGlyphOffsets = NULL;

NK_API nk_handle nk_sdl_get_userdata(struct nk_context *ctx)
{
   struct nk_sdl *sdl;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);
   return sdl->userdata;
}

NK_API void nk_sdl_set_userdata(struct nk_context *ctx, nk_handle userdata)
{
   struct nk_sdl *sdl;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
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
   allocator.userdata.ptr = 0;
   allocator.alloc = nk_sdl_alloc;
   allocator.free = nk_sdl_free;
   return allocator;
}

NK_INTERN void nk_sdl_device_upload_atlas(struct nk_context *ctx, const void *image, int width, int height)
{
   struct nk_sdl *sdl;
   NK_ASSERT(ctx);
   NK_ASSERT(image);
   NK_ASSERT(width > 0);
   NK_ASSERT(height > 0);

   sdl = (struct nk_sdl *)ctx->userdata.ptr;
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
   struct nk_sdl *sdl;
   bool active;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   /* Determine if Nuklear is using any top-level "edit" widget.
    * Popups take higher priority because they block any incomming input.
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
      const bool window_edit_active = SDL_TextInputActive(sdl->win);

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
   /* setup global state */
   struct nk_sdl *sdl;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   { /* setup internal delta time that Nuklear needs for animations */
      const Uint64 now = SDL_GetTicksNS();
      ctx->delta_time_seconds = (float)(now - sdl->last_render) / (float)SDL_NS_PER_SECOND;
      sdl->last_render = now;
   }

   {
      SDL_Rect saved_clip;
      bool clipping_enabled;
      int vs = sizeof(struct nk_sdl_vertex);
      size_t vp = NK_OFFSETOF(struct nk_sdl_vertex, position);
      size_t vt = NK_OFFSETOF(struct nk_sdl_vertex, uv);
      size_t vc = NK_OFFSETOF(struct nk_sdl_vertex, col);

      /* convert from command queue into draw list and draw to screen */
      const struct nk_draw_command *cmd;
      const nk_draw_index *offset = NULL;
      struct nk_buffer vbuf, ebuf;

      /* fill converting configuration */
      struct nk_convert_config config;
      static const struct nk_draw_vertex_layout_element vertex_layout[] = {{NK_VERTEX_POSITION, NK_FORMAT_FLOAT, NK_OFFSETOF(struct nk_sdl_vertex, position)},
         {NK_VERTEX_TEXCOORD, NK_FORMAT_FLOAT, NK_OFFSETOF(struct nk_sdl_vertex, uv)},
         {NK_VERTEX_COLOR, NK_FORMAT_R32G32B32A32_FLOAT, NK_OFFSETOF(struct nk_sdl_vertex, col)}, {NK_VERTEX_LAYOUT_END}};
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
      nk_buffer_init(&vbuf, &sdl->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
      nk_buffer_init(&ebuf, &sdl->allocator, NK_BUFFER_DEFAULT_INITIAL_SIZE);
      nk_convert(&sdl->ctx, &sdl->ogl.cmds, &vbuf, &ebuf, &config);

      /* iterate over and execute each draw command */
      offset = (const nk_draw_index *)nk_buffer_memory_const(&ebuf);

      clipping_enabled = SDL_RenderClipEnabled(sdl->renderer);
      SDL_GetRenderClipRect(sdl->renderer, &saved_clip);

      nk_draw_foreach(cmd, &sdl->ctx, &sdl->ogl.cmds)
      {
         if (!cmd->elem_count)
            continue;

         {
            SDL_Rect r;
            r.x = cmd->clip_rect.x;
            r.y = cmd->clip_rect.y;
            r.w = cmd->clip_rect.w;
            r.h = cmd->clip_rect.h;
            SDL_SetRenderClipRect(sdl->renderer, &r);
         }

         {
            const void *vertices = nk_buffer_memory_const(&vbuf);

            SDL_RenderGeometryRaw(sdl->renderer, (SDL_Texture *)cmd->texture.ptr, (const float *)((const nk_byte *)vertices + vp), vs,
               (const SDL_FColor *)((const nk_byte *)vertices + vc), vs, (const float *)((const nk_byte *)vertices + vt), vs, (vbuf.needed / vs),
               (void *)offset, cmd->elem_count, 2);

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
   char *text;
   int len;
   NK_UNUSED(usr);

   /* this function returns empty string on failure, not NULL */
   text = SDL_GetClipboardText();
   NK_ASSERT(text);

   if (text[0] != '\0')
   {
      /* FIXME: there is a bug in Nuklear that affects UTF8 clipboard handling
       * "len" should be a buffer length, but due to bug it must be a glyph count
       * see: https://github.com/Immediate-Mode-UI/Nuklear/pull/841 */
#if 0
        len = nk_strlen(text);
#else
      len = SDL_utf8strlen(text);
#endif
      nk_textedit_paste(edit, text, len);
   }
   SDL_free(text);
}

NK_INTERN void nk_sdl_clipboard_copy(nk_handle usr, const char *text, int len)
{
   const char *ptext;
   char *str;
   size_t buflen;
   int i;
   struct nk_sdl *sdl = (struct nk_sdl *)usr.ptr;
   NK_ASSERT(sdl);
   if (len <= 0 || text == NULL)
      return;

   /* FIXME: there is a bug in Nuklear that affects UTF8 clipboard handling
    * "len" is expected to be a buffer length, but due to bug it actually is a glyph count
    * see: https://github.com/Immediate-Mode-UI/Nuklear/pull/841 */
#if 0
    buflen = len + 1;
    NK_UNUSED(ptext);
#else
   ptext = text;
   for (i = len; i > 0; i--)
      (void)SDL_StepUTF8(&ptext, NULL);
   buflen = (size_t)(ptext - text) + 1;
#endif

   str = (char *)sdl->allocator.alloc(sdl->allocator.userdata, 0, buflen);
   if (!str)
      return;
   SDL_strlcpy(str, text, buflen);
   SDL_SetClipboardText(str);
   sdl->allocator.free(sdl->allocator.userdata, str);
}

NK_API struct nk_context *nk_sdl_init(SDL_Window *win, SDL_Renderer *renderer, struct nk_allocator allocator)
{
   struct nk_sdl *sdl;
   NK_ASSERT(win);
   NK_ASSERT(renderer);
   NK_ASSERT(allocator.alloc);
   NK_ASSERT(allocator.free);
   sdl = (struct nk_sdl *)allocator.alloc(allocator.userdata, 0, sizeof(*sdl));
   NK_ASSERT(sdl);
   SDL_zerop(sdl);
   sdl->allocator.userdata = allocator.userdata;
   sdl->allocator.alloc = allocator.alloc;
   sdl->allocator.free = allocator.free;
   sdl->win = win;
   sdl->renderer = renderer;
   nk_init(&sdl->ctx, &sdl->allocator, 0);
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

#ifdef NK_INCLUDE_FONT_BAKING
NK_API struct nk_font_atlas *nk_sdl_font_stash_begin(struct nk_context *ctx)
{
   struct nk_sdl *sdl;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);
   nk_font_atlas_init(&sdl->atlas, &sdl->allocator);
   nk_font_atlas_begin(&sdl->atlas);
   return &sdl->atlas;
}
#endif

#ifdef NK_INCLUDE_FONT_BAKING
NK_API void nk_sdl_font_stash_end(struct nk_context *ctx)
{
   struct nk_sdl *sdl;
   const void *image;
   int w, h;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);
   image = nk_font_atlas_bake(&sdl->atlas, &w, &h, NK_FONT_ATLAS_RGBA32);
   NK_ASSERT(image);
   nk_sdl_device_upload_atlas(&sdl->ctx, image, w, h);
   nk_font_atlas_end(&sdl->atlas, nk_handle_ptr(sdl->ogl.font_tex), &sdl->ogl.tex_null);
   if (sdl->atlas.default_font)
   {
      nk_style_set_font(&sdl->ctx, &sdl->atlas.default_font->handle);
   }
}
#endif

NK_API int nk_sdl_handle_event(struct nk_context *ctx, SDL_Event *evt)
{
   struct nk_sdl *sdl;

   NK_ASSERT(ctx);
   NK_ASSERT(evt);

   sdl = (struct nk_sdl *)ctx->userdata.ptr;
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
      int down = evt->type == SDL_EVENT_KEY_DOWN;
      int ctrl_down = evt->key.mod & SDL_KMOD_CTRL;

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
      const int x = evt->button.x, y = evt->button.y;
      const int down = evt->button.down;
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
      nk_glyph glyph;
      nk_size len;
      NK_ASSERT(evt->text.text);
      len = SDL_strlen(evt->text.text);
      NK_ASSERT(len <= NK_UTF_SIZE);
      NK_MEMCPY(glyph, evt->text.text, len);
      nk_input_glyph(ctx, glyph);
   }
      return 1;

   case SDL_EVENT_MOUSE_WHEEL:
      nk_input_scroll(ctx, nk_vec2(evt->wheel.x, evt->wheel.y));
      return 1;
   }
   return 0;
}

NK_API
void nk_sdl_shutdown(struct nk_context *ctx)
{
   struct nk_sdl *sdl;
   NK_ASSERT(ctx);
   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

#ifdef NK_INCLUDE_FONT_BAKING
   if (sdl->atlas.font_num > 0)
      nk_font_atlas_clear(&sdl->atlas);
#endif

   nk_buffer_free(&sdl->ogl.cmds);

   if (sdl->ogl.font_tex != NULL)
   {
      SDL_DestroyTexture(sdl->ogl.font_tex);
      sdl->ogl.font_tex = NULL;
   }

   nk_free(ctx);
   sdl->allocator.free(sdl->allocator.userdata, sdl->debug_font);
   sdl->allocator.free(sdl->allocator.userdata, sdl);
}

/* Debug Font Width/Height of internal texture atlas
 * This is a result of: ceil(sqrt('~' - ' '))
 * There is a sanity check for this value in nk_sdl_style_set_debug_font */
#define NK_SDL_DFWH (10)

NK_INTERN float nk_sdl_query_debug_font_width(nk_handle handle, float height, const char *text, int len)
{
   NK_UNUSED(handle);
   return nk_utf_len(text, len) * (height + 1) + 1; // one per character for border (which may overlap) plus one for last character non-overlap.
}

NK_INTERN void nk_sdl_query_debug_font_glypth(nk_handle handle, float height, struct nk_user_font_glyph *glyph, nk_rune codepoint, nk_rune next_codepoint)
{
   char ascii;
   int idx, x, y;
   NK_UNUSED(next_codepoint);
   NK_UNUSED(handle);

   /* replace non-ASCII characters with question mark */
   ascii = (codepoint < (nk_rune)' ' || codepoint > (nk_rune)'~') ? '?' : (char)codepoint;
   NK_ASSERT(ascii >= ' ' && ascii <= '~');

   idx = (int)(ascii - ' ');
   x = idx / NK_SDL_DFWH;
   y = idx % NK_SDL_DFWH;
   NK_ASSERT(x >= 0 && x < NK_SDL_DFWH);
   NK_ASSERT(y >= 0 && y < NK_SDL_DFWH);

   glyph->height = height;
   glyph->width = height;
   glyph->xadvance = height - 2; // TODO: Is this shift by 2 correct? probably, because border is allowed to overlap.
   glyph->uv[0].x = (float)(x + 0) / NK_SDL_DFWH;
   glyph->uv[0].y = (float)(y + 0) / NK_SDL_DFWH;
   glyph->uv[1].x = (float)(x + 1) / NK_SDL_DFWH;
   glyph->uv[1].y = (float)(y + 1) / NK_SDL_DFWH;
   glyph->offset.x = 0.0f;
   glyph->offset.y = 0.0f;
}

NK_API void nk_sdl_style_set_debug_font(struct nk_context *ctx)
{
   struct nk_user_font *font;
   struct nk_sdl *sdl;
   SDL_Surface *surface;
   SDL_Renderer *renderer;
   char buf[2];
   int x, y;
   bool success;
   NK_ASSERT(ctx);

   sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   if (sdl->debug_font)
   {
      sdl->allocator.free(sdl->allocator.userdata, sdl->debug_font);
      sdl->debug_font = 0;
   }

   /* sanity check: formal proof of NK_SDL_DFWH value (which is 10) */
   NK_ASSERT(SDL_ceil(SDL_sqrt('~' - ' ')) == NK_SDL_DFWH);

   /* We use another Software Renderer just to make sure
    * that we won't mutate any state in the main Renderer. */
   // 4 = 2 for border, plus 2 for buffer to avoid blur bleed
   surface = SDL_CreateSurface(
      NK_SDL_DFWH * (SDL_DEBUG_TEXT_FONT_CHARACTER_SIZE + 4), NK_SDL_DFWH * (SDL_DEBUG_TEXT_FONT_CHARACTER_SIZE + 4), SDL_PIXELFORMAT_RGBA32);
   NK_ASSERT(surface);
   renderer = SDL_CreateSoftwareRenderer(surface);
   NK_ASSERT(renderer);
   // success = SDL_SetRenderDrawColor(renderer, 0xFF, 0xFF, 0xFF, 0xFF);
   // NK_ASSERT(success);

   /* SPACE is the first printable ASCII character */
   NK_MEMCPY(buf, " ", sizeof(buf));
   for (x = 0; x < NK_SDL_DFWH; x++)
   {
      for (y = 0; y < NK_SDL_DFWH; y++)
      {
         int startX = x * (SDL_DEBUG_TEXT_FONT_CHARACTER_SIZE + 4) + 2;
         int startY = y * (SDL_DEBUG_TEXT_FONT_CHARACTER_SIZE + 4) + 2;
         success = SDL_SetRenderDrawColor(renderer, 0x00, 0x00, 0x00, 0xFF);
         for (int i = 0; i < 9; i++)
         {
            int offX = (i % 3) - 1;
            int offY = (i / 3) - 1;
            if (offX != 0 || offY != 0)
            {
               SDL_RenderDebugText(renderer, (float)(startX + offX), (float)(startY + offY), buf);
            }
         }
         success = SDL_SetRenderDrawColor(renderer, 0x1E, 0xAC, 0x5B, 0xFF); // TODO: Why is this BGRA and not RGBA?
         success = SDL_RenderDebugText(renderer, (float)(startX), (float)(startY), buf);
         NK_ASSERT(success);
         buf[0]++;

         /* TILDE is the last printable ASCII character */
         if (buf[0] > '~')
            break;
      }
   }
   success = SDL_RenderPresent(renderer);
   NK_ASSERT(success);

   font = (struct nk_user_font *)sdl->allocator.alloc(sdl->allocator.userdata, 0, sizeof(*font));
   NK_ASSERT(font);
   font->userdata.ptr = sdl;
   font->height = SDL_DEBUG_TEXT_FONT_CHARACTER_SIZE;
   font->width = &nk_sdl_query_debug_font_width;
   font->query = &nk_sdl_query_debug_font_glypth;

   /* HACK: nk_sdl_device_upload_atlas turns pixels into SDL_Texture
    *       and sets said Texture into sdl->ogl.font_tex
    *       then nk_sdl_render expects same Texture at font->texture */
   nk_sdl_device_upload_atlas(ctx, surface->pixels, surface->w, surface->h);
   font->texture.ptr = sdl->ogl.font_tex;

   sdl->debug_font = font;
   nk_style_set_font(ctx, font);

   SDL_DestroyRenderer(renderer);
   SDL_DestroySurface(surface);
}

NK_INTERN float nk_sdl_query_tiny_font_width(nk_handle handle, float height, char const *text, int len)
{
   NK_UNUSED(handle);
   int32_t width = 0;
   char const *const end = text + len;
   for (char const *it = text; it != end; ++it)
   {
      // TODO this is wrong, as the text is UTF-8, and it needs to be decoded and then translated into the right codepage.
      char const ascii = (*it < tinyFontCharLowest || *it > (nk_rune)tinyFontCharHighest) ? '?' : *it;

      size_t const characterIndex = ascii - tinyFontCharLowest;

      // TODO: try to figure out why -1 makes it work suddenly with a text input (well, almost; new lines mess it up)
      width += (tinyFontGlyphOffsets[characterIndex + 1] - tinyFontGlyphOffsets[characterIndex]) - 1;
   }
   float scale = height / (float)(tinyFontHeight + 2);
   return (float)(width + 1) * scale;
}

NK_INTERN void nk_sdl_query_tiny_font_glyph(nk_handle handle, float height, struct nk_user_font_glyph *glyph, nk_rune codepoint, nk_rune next_codepoint)
{
   NK_UNUSED(next_codepoint);
   NK_UNUSED(handle);

   // replace unknown characters with question mark
   char const ascii = (codepoint < (nk_rune)tinyFontCharLowest || codepoint > (nk_rune)tinyFontCharHighest) ? '?' : (char)codepoint;

   int const characterCount = tinyFontCharHighest - tinyFontCharLowest + 1;

   // +2 per character for border, +2 per character for buffer to avoid blur bleed
   int const bitmapWidth = (int)(tinyFontGlyphOffsets[characterCount + 1] - tinyFontGlyphOffsets[0]) + (characterCount * 2) + (characterCount * 2);

   size_t const characterIndex = ascii - tinyFontCharLowest;

   ptrdiff_t const glyphOffsetBegin = tinyFontGlyphOffsets[characterIndex];
   ptrdiff_t const glyphOffsetEnd = tinyFontGlyphOffsets[characterIndex + 1];
   int32_t width = glyphOffsetEnd - glyphOffsetBegin;

   float scale = height / (float)(tinyFontHeight + 2);
   glyph->height = (float)(tinyFontHeight + 2) * scale;
   glyph->width = (float)(width + 1) * scale;
   glyph->xadvance = (float)width * scale;
   glyph->uv[0].x = (float)(glyphOffsetBegin + (characterIndex * 2) + (characterIndex * 2) + 1) / (float)bitmapWidth;
   glyph->uv[0].y = 0;
   glyph->uv[1].x = (float)(glyphOffsetEnd + (characterIndex * 2) + (characterIndex * 2) + 2) / (float)bitmapWidth;
   glyph->uv[1].y = 1;
   glyph->offset.x = 0.0f;
   glyph->offset.y = 0.0f;
}

static uint16_t serialReadU16LittleEndian(void const *addr)
{
   uint8_t const *data = addr;
   return ((uint16_t)data[1] << 8) | (uint16_t)data[0];
}

static int16_t serialReadS16LittleEndian(void const *addr)
{
   return (int16_t)serialReadU16LittleEndian(addr);
}

static uint32_t serialReadU32LittleEndian(void const *addr)
{
   uint8_t const *data = addr;
   return ((uint32_t)data[3] << 24) | ((uint32_t)data[2] << 16) | ((uint32_t)data[1] << 8) | (uint32_t)data[0];
}

static int32_t serialReadS32LittleEndian(void const *addr)
{
   return (int32_t)serialReadU32LittleEndian(addr);
}

typedef int16_t PixelAxisPosition;
typedef int32_t PixelAxisOffset;

typedef struct
{
   PixelAxisPosition x;
   PixelAxisPosition y;
} PixelPosition;

typedef struct
{
   PixelAxisOffset x;
   PixelAxisOffset y;
} PixelOffset;

typedef struct
{
   uint8_t const *data;
   ptrdiff_t stride;
} Bitmap;

typedef struct
{
   uint8_t r;
   uint8_t g;
   uint8_t b;
} ColorRgb;

static void renderMonochromeFontCharacter(SDL_Surface *surface, size_t const characterIndex, PixelOffset const off, Bitmap const bitmap, ColorRgb const color)
{
   // TODO figure out types; when which is more appropriate.
   ptrdiff_t const glyphOffsetBegin = tinyFontGlyphOffsets[characterIndex];
   ptrdiff_t const glyphOffsetEnd = tinyFontGlyphOffsets[characterIndex + 1];
   int16_t const width = (int16_t)(glyphOffsetEnd - glyphOffsetBegin);
   int outXBase = glyphOffsetBegin + (characterIndex * 2) + (characterIndex * 2) + 2 + off.x;
   int outYBase = 1 + off.y;
   for (int16_t y = 0; y < tinyFontHeight; ++y)
   {
      for (int16_t x = 0; x < width; ++x)
      {
         size_t const byteOffset = (glyphOffsetBegin + x) / 8;
         uint8_t const in = *(bitmap.data + (bitmap.stride * y) + byteOffset);
         size_t const bitOffset = (glyphOffsetBegin + x) % 8;
         if ((in & (0x80 >> bitOffset)) != 0)
         {
            SDL_WriteSurfacePixel(surface, outXBase + x, outYBase + y, color.b, color.g, color.r, 0xFF); // TODO: Why is this BGRA and not RGBA?
         }
      }
   }
}

NK_API void nk_sdl_style_set_tiny_font(struct nk_context *ctx, float scale)
{
   NK_ASSERT(ctx);

   struct nk_sdl *sdl = (struct nk_sdl *)ctx->userdata.ptr;
   NK_ASSERT(sdl);

   if (sdl->debug_font)
   {
      sdl->allocator.free(sdl->allocator.userdata, sdl->debug_font);
      sdl->debug_font = 0;
   }

   tinyFontCharLowest = serialReadS16LittleEndian(tinyFont + 0x0024);
   tinyFontCharHighest = serialReadS16LittleEndian(tinyFont + 0x0026);
   ptrdiff_t const tinyFontOffset = serialReadS32LittleEndian(tinyFont + 0x0048);
   tinyFontGlyphOffsets = (int16_t *)(tinyFont + tinyFontOffset);
   ptrdiff_t const tinyFontBitmapOffset = serialReadS32LittleEndian(tinyFont + 0x004C);
   uint8_t const *tinyFontBitmap = tinyFont + tinyFontBitmapOffset;
   int16_t const tinyFontBitmapWidth = serialReadS16LittleEndian(tinyFont + 0x0050);
   tinyFontHeight = serialReadS16LittleEndian(tinyFont + 0x0052);
   int const characterCount = tinyFontCharHighest - tinyFontCharLowest + 1;

   // +2 per character for border, +2 per character for buffer to avoid blur bleed
   int bitmapWidth = (int)(tinyFontGlyphOffsets[characterCount + 1] - tinyFontGlyphOffsets[0]) + (characterCount * 2) + (characterCount * 2);
   int bitmapHeight = tinyFontHeight + 2;
   SDL_Surface *surface = SDL_CreateSurface(bitmapWidth, bitmapHeight, SDL_PIXELFORMAT_RGBA32);
   NK_ASSERT(surface);

   Bitmap fontBitmap = {.data = tinyFontBitmap, .stride = tinyFontBitmapWidth};
   ColorRgb black = {};

   for (size_t currentCharOffset = 0; currentCharOffset < characterCount; ++currentCharOffset)
   {
      for (PixelAxisOffset i = 0; i < 9; ++i)
      {
         PixelOffset const off = {.x = (PixelAxisOffset)((i % 3) - 1), .y = (PixelAxisOffset)((i / 3) - 1)};
         if (off.x != 0 || off.y != 0)
         {
            renderMonochromeFontCharacter(surface, currentCharOffset, off, fontBitmap, black);
         }
      }
      ColorRgb textColor = {.r = 0x5B, .g = 0xAC, .b = 0x1E};
      PixelOffset zeroOffset = {};
      renderMonochromeFontCharacter(surface, currentCharOffset, zeroOffset, fontBitmap, textColor);
   }

   struct nk_user_font *font = (struct nk_user_font *)sdl->allocator.alloc(sdl->allocator.userdata, NULL, sizeof(*font));
   NK_ASSERT(font);
   font->userdata.ptr = sdl;
   font->height = (float)(tinyFontHeight + 2) * scale;
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
