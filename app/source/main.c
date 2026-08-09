#define SDL_MAIN_USE_CALLBACKS 1
#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

#include "dcimgui.h"
#include "dcimgui_impl_sdl3.h"
#include "dcimgui_impl_sdlrenderer3.h"

struct SDLApp
{
   SDL_Window *window;
   SDL_Renderer *renderer;
};

static SDL_AppResult sdlFail(char const *const message)
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
      return sdlFail("failed to initialize SDL");
   }

   SDL_SetAppMetadata("InkyBlackness - HackEd", REPO_SHORT_VERSION, "io.github.inkyblackness.hacked");
   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return sdlFail("failed to initialize SDL");
   }
   struct SDLApp *app = SDL_malloc(sizeof(*app));
   if (app == NULL)
   {
      return sdlFail("failed to allocate application memory");
   }
   if (!SDL_CreateWindowAndRenderer(
          "InkyBlackness - HackEd - " REPO_LONG_VERSION, 320, 200, SDL_WINDOW_RESIZABLE | SDL_WINDOW_HIGH_PIXEL_DENSITY, &app->window, &app->renderer))
   {
      SDL_free(app);
      return sdlFail("failed to create window/renderer");
   }
   *appstate = app;

   if (!SDL_SetRenderVSync(app->renderer, 1))
   {
      SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "SDL_SetRenderVSync failed: %s", SDL_GetError());
   }
   SDL_SetRenderLogicalPresentation(app->renderer, 320, 200, SDL_LOGICAL_PRESENTATION_LETTERBOX);

   CIMGUI_CHECKVERSION();
   ImGui_CreateContext(NULL);
   ImGuiIO *io = ImGui_GetIO();
   io->ConfigFlags |= ImGuiConfigFlags_NavEnableKeyboard;
   io->ConfigFlags |= ImGuiConfigFlags_DockingEnable;
   io->IniFilename = NULL;
   ImGui_StyleColorsDark(NULL);

   cImGui_ImplSDL3_InitForSDLRenderer(app->window, app->renderer);
   cImGui_ImplSDLRenderer3_Init(app->renderer);

   return SDL_APP_CONTINUE;
}

SDL_AppResult SDL_AppEvent(void *appstate, SDL_Event *event)
{
   struct SDLApp *app = (struct SDLApp *)appstate;

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
      // You may wish to rescale the renderer.
      return SDL_APP_CONTINUE;
   default:
      break;
   }
   // Remember to always rescale the event coordinates, if your renderer uses custom scale.
   SDL_ConvertEventToRenderCoordinates(app->renderer, event);

   cImGui_ImplSDL3_ProcessEvent(event);

   return SDL_APP_CONTINUE;
}

static void appClearBackground(SDL_Renderer *const renderer)
{
   SDL_SetRenderDrawColor(renderer, 0x00, 0x00, 0x00, 0xFF);
   SDL_RenderClear(renderer);
}

SDL_AppResult SDL_AppIterate(void *appstate)
{
   struct SDLApp *app = (struct SDLApp *)appstate;
   SDL_AppResult appResult = SDL_APP_CONTINUE;
   // float const scale = SDL_GetWindowDisplayScale(app->window);

   {
      cImGui_ImplSDLRenderer3_NewFrame();
      cImGui_ImplSDL3_NewFrame();
      ImGui_NewFrame();

      bool showDemoWindow = true;
      ImGui_ShowDemoWindow(&showDemoWindow);

      ImGui_Render();
   }

   appClearBackground(app->renderer);

   // SDL_SetRenderScale(app->renderer, io.DisplayFramebufferScale.x, io.DisplayFramebufferScale.y); // TODO: needed?
   cImGui_ImplSDLRenderer3_RenderDrawData(ImGui_GetDrawData(), app->renderer);
   SDL_RenderPresent(app->renderer);

   return appResult;
}

void SDL_AppQuit(void *appstate, SDL_AppResult result)
{
   (void)result;

   struct SDLApp *app = (struct SDLApp *)appstate;

   if (app)
   {
      SDL_Log("Quitting");

      cImGui_ImplSDLRenderer3_Shutdown();
      cImGui_ImplSDL3_Shutdown();
      ImGui_DestroyContext(ImGui_GetCurrentContext());

      SDL_DestroyRenderer(app->renderer);
      SDL_DestroyWindow(app->window);
      SDL_free(app);
   }
}
