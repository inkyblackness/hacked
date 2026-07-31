#pragma once

#include "hacked/nuklear/Include.h"

#ifdef __cplusplus
extern "C" {
#endif

NK_API struct nk_context *nk_sdl_init(SDL_Window *win, SDL_Renderer *renderer, struct nk_allocator allocator);
#ifdef NK_INCLUDE_FONT_BAKING
NK_API struct nk_font_atlas *nk_sdl_font_stash_begin(struct nk_context *ctx);
NK_API void nk_sdl_font_stash_end(struct nk_context *ctx);
#endif
NK_API int nk_sdl_handle_event(struct nk_context *ctx, SDL_Event *evt);
NK_API void nk_sdl_render(struct nk_context *ctx, enum nk_anti_aliasing);
NK_API void nk_sdl_update_TextInput(struct nk_context *ctx);
NK_API void nk_sdl_shutdown(struct nk_context *ctx);
NK_API nk_handle nk_sdl_get_userdata(struct nk_context *ctx);
NK_API void nk_sdl_set_userdata(struct nk_context *ctx, nk_handle userdata);
NK_API void nk_sdl_style_set_debug_font(struct nk_context *ctx);
NK_API struct nk_allocator nk_sdl_allocator(void);

#ifdef __cplusplus
}
#endif
