#pragma once

#include "hacked/nuklear/Include.h"

#ifdef __cplusplus
extern "C" {
#endif

NK_API struct nk_context *uiBridgeInit(SDL_Window *win, SDL_Renderer *renderer);
NK_API void uiBridgeShutdown(struct nk_context *ctx);
NK_API nk_handle uiBridgeGetUserdata(struct nk_context const *ctx);
NK_API void uiBridgeSetUserdata(struct nk_context *ctx, nk_handle userdata);
NK_API void uiBridgeSetFont(struct nk_context *ctx, float scale);

NK_API int uiBridgeHandleEvent(struct nk_context *ctx, SDL_Event const *evt);
NK_API void uiBridgeUpdateTextInput(struct nk_context *ctx);

NK_API void uiBridgeRender(struct nk_context *ctx);

#ifdef __cplusplus
}
#endif
