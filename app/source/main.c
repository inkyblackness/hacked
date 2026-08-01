#define SDL_MAIN_USE_CALLBACKS 1
#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

#include "hacked/nuklear/NuklearSdlBridge.h"

struct nk_sdl_app
{
   SDL_Window *window;
   SDL_Renderer *renderer;
   struct nk_context *ctx;
   enum nk_anti_aliasing AA;
};

static SDL_AppResult nk_sdl_fail(char const *const message)
{
   SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "Error; %s: %s", message, SDL_GetError());
   return SDL_APP_FAILURE;
}

SDL_AppResult SDL_AppInit(void **const appstate, int const argc, char *argv[])
{
   (void)argc;
   (void)argv;

   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return nk_sdl_fail("failed to initialize SDL");
   }

   SDL_SetAppMetadata("InkyBlackness - HackEd", REPO_SHORT_VERSION, "io.github.inkyblackness.hacked");
   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return nk_sdl_fail("failed to initialize SDL");
   }
   struct nk_sdl_app *app = SDL_malloc(sizeof(*app));
   if (app == NULL)
   {
      return nk_sdl_fail("failed to allocate application memory");
   }
   if (!SDL_CreateWindowAndRenderer("InkyBlackness - HackEd - " REPO_LONG_VERSION, 320, 200, SDL_WINDOW_RESIZABLE, &app->window, &app->renderer))
   {
      SDL_free(app);
      return nk_sdl_fail("failed to create window/renderer");
   }
   *appstate = app;

   if (!SDL_SetRenderVSync(app->renderer, 1))
   {
      SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "SDL_SetRenderVSync failed: %s", SDL_GetError());
   }
   SDL_SetRenderLogicalPresentation(app->renderer, 320, 200, SDL_LOGICAL_PRESENTATION_LETTERBOX);

   struct nk_context *ctx = nk_sdl_init(app->window, app->renderer, nk_sdl_allocator());
   app->ctx = ctx;

   {
      nk_sdl_style_set_tiny_font(ctx, 1.0f);

      /* You may wish to change a few style options, here are few recommendations: */
      ctx->style.button.rounding = 0.0f;
      ctx->style.menu_button.rounding = 0.0f;
      ctx->style.property.rounding = 0.0f;
      ctx->style.property.border = 0.0f;
      ctx->style.option.border = -1.0f;
      ctx->style.checkbox.border = -1.0f;
      ctx->style.property.dec_button.border = -2.0f;
      ctx->style.property.inc_button.border = -2.0f;
      ctx->style.tab.tab_minimize_button.border = -2.0f;
      ctx->style.tab.tab_maximize_button.border = -2.0f;
      ctx->style.tab.node_minimize_button.border = -2.0f;
      ctx->style.tab.node_maximize_button.border = -2.0f;
      ctx->style.checkbox.spacing = 5.0f;
      ctx->style.window.header.label_padding.x = 0.0f;
      ctx->style.window.header.label_padding.y = 0.0f;
      ctx->style.window.header.padding.x = 0.0f;
      ctx->style.window.header.padding.y = 0.0f;
      ctx->style.window.header.spacing.x = 0.0f;
      ctx->style.window.header.spacing.y = 0.0f;

      /* It's better to disable anti-aliasing when using small fonts */
      app->AA = NK_ANTI_ALIASING_OFF;
   }
   nk_input_begin(ctx);

   return SDL_APP_CONTINUE;
}

SDL_AppResult SDL_AppEvent(void *appstate, SDL_Event *event)
{
   struct nk_sdl_app *app = (struct nk_sdl_app *)appstate;

   switch (event->type)
   {
   case SDL_EVENT_QUIT:
      return SDL_APP_SUCCESS;
   case SDL_EVENT_KEY_DOWN:
      if (event->key.key == SDLK_Q && event->key.mod & SDL_KMOD_CTRL)
      {
         return SDL_APP_SUCCESS;
      }
      break;
   case SDL_EVENT_WINDOW_DISPLAY_SCALE_CHANGED:
      /* You may wish to rescale the renderer and Nuklear during this event.
       * Without this the UI and Font could appear too small or too big.
       * This is not handled by the demo in order to keep it simple,
       * but you may wish to re-bake the Font whenever this happens. */
      SDL_Log("Unhandled scale event! Nuklear may appear blurry");
      return SDL_APP_CONTINUE;
   default:
      break;
   }

   /* Remember to always rescale the event coordinates,
    * if your renderer uses custom scale. */
   SDL_ConvertEventToRenderCoordinates(app->renderer, event);

   nk_sdl_handle_event(app->ctx, event);

   return SDL_APP_CONTINUE;
}

SDL_AppResult SDL_AppIterate(void *appstate)
{
   struct nk_sdl_app *app = (struct nk_sdl_app *)appstate;
   struct nk_context *ctx = app->ctx;

   const double now = ((double)SDL_GetTicks()) / 1000.0;
   /* choose the color for the frame we will draw. The sine wave trick makes it fade between colors smoothly. */
   const float red = (float)(0.5 + 0.5 * SDL_sin(now));
   const float green = (float)(0.5 + 0.5 * SDL_sin(now + SDL_PI_D * 2 / 3));
   const float blue = (float)(0.5 + 0.5 * SDL_sin(now + SDL_PI_D * 4 / 3));

   nk_input_end(ctx);
   /* GUI */
   if (nk_begin(ctx, "Demo", nk_rect(50, 50, 230, 190), NK_WINDOW_BORDER | NK_WINDOW_MOVABLE | NK_WINDOW_SCALABLE | NK_WINDOW_MINIMIZABLE | NK_WINDOW_TITLE))
   {
      enum
      {
         EASY,
         HARD
      };
      static int op = EASY;
      static int property = 20;
      enum
      {
         textBufferSize = 1000
      };
      static char textBuffer[textBufferSize] = "";
      static int textBufferUsedLen = 0;

      nk_layout_row_static(ctx, 9, 60, 1);
      if (nk_button_label(ctx, "button"))
      {
         SDL_Log("button pressed");
      }
      nk_layout_row_dynamic(ctx, 30, 2);
      if (nk_option_label(ctx, "easy", op == EASY))
         op = EASY;
      if (nk_option_label(ctx, "hard", op == HARD))
         op = HARD;
      nk_layout_row_dynamic(ctx, 25, 1);
      nk_property_int(ctx, "Compression:", 0, &property, 1000, 1, 1);

      nk_layout_row_dynamic(ctx, 20, 1);
      nk_label(ctx, "background:", NK_TEXT_LEFT);
      nk_layout_row_dynamic(ctx, 25, 1);

      nk_edit_string(ctx, NK_EDIT_ALWAYS_INSERT_MODE | NK_EDIT_MULTILINE, textBuffer, &textBufferUsedLen, textBufferSize, NULL);
   }
   nk_end(ctx);

   SDL_SetRenderDrawColorFloat(app->renderer, red, green, blue, SDL_ALPHA_OPAQUE_FLOAT); /* new color, full alpha. */
   SDL_RenderClear(app->renderer);

   nk_sdl_render(ctx, app->AA);
   nk_sdl_update_TextInput(ctx);

   /* show if TextInput is active for debug purpose. Feel free to remove this. */
   SDL_SetRenderDrawColor(app->renderer, 0xFF, 0xFF, 0xFF, 0xFF);
   SDL_RenderDebugTextFormat(app->renderer, 10, 10, "TextInputActive? %s", SDL_TextInputActive(app->window) ? "Yes" : "No");

   SDL_RenderPresent(app->renderer);

   nk_input_begin(ctx);
   return SDL_APP_CONTINUE;
}

void SDL_AppQuit(void *appstate, SDL_AppResult result)
{
   (void)result;

   struct nk_sdl_app *app = (struct nk_sdl_app *)appstate;

   if (app)
   {
      SDL_Log("Quitting");
      nk_input_end(app->ctx);
      nk_sdl_shutdown(app->ctx);
      SDL_DestroyRenderer(app->renderer);
      SDL_DestroyWindow(app->window);
      SDL_free(app);
   }
}

char *nk_sdl_dtoa(char *str, double d)
{
   NK_ASSERT(str);
   if (!str)
      return NULL;
   (void)SDL_snprintf(str, 99999, "%.17g", d);
   return str;
}
