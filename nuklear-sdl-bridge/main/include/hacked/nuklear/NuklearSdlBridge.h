#pragma once

#include "hacked/nuklear/Include.h"

#ifdef __cplusplus
extern "C" {
#endif

NK_API struct nk_context *nk_sdl_init(SDL_Window *win, SDL_Renderer *renderer, struct nk_allocator allocator);
NK_API int nk_sdl_handle_event(struct nk_context *ctx, SDL_Event *evt);
NK_API void nk_sdl_render(struct nk_context *ctx);
NK_API void nk_sdl_update_TextInput(struct nk_context *ctx);
NK_API void nk_sdl_shutdown(struct nk_context *ctx);
NK_API nk_handle nk_sdl_get_userdata(struct nk_context *ctx);
NK_API void nk_sdl_set_userdata(struct nk_context *ctx, nk_handle userdata);
NK_API void nk_sdl_style_set_tiny_font(struct nk_context *ctx, float scale);
NK_API struct nk_allocator nk_sdl_allocator();

#ifdef __cplusplus
}
#endif
