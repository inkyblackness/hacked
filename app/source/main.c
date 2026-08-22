#include <float.h>
#include <stdint.h>
#include <stdio.h>

#define SDL_MAIN_USE_CALLBACKS 1
#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

#include "dcimgui.h"
#include "dcimgui_impl_sdl3.h"
#include "dcimgui_impl_sdlrenderer3.h"

struct String
{
   char *text;
   size_t len;
};

struct HackEdApp
{
   int originalDesktopWidth;
   int originalDesktopHeight;

   SDL_Window *window;
   SDL_Renderer *renderer;

   struct String folderHome;
   struct String folderDocuments;
   struct String folderApp;
   struct String folderPreferences;

   bool showDemoWindow;
};

static float appGetBaseUIScale(SDL_Window *const window)
{
   return SDL_GetWindowDisplayScale(window);
}

static SDL_AppResult appFailSDL(char const *const message)
{
   SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "Error; %s: %s", message, SDL_GetError());
   return SDL_APP_FAILURE;
}

static ImVec4 color(uint8_t const r, uint8_t const g, uint8_t const b, float const a)
{
   ImVec4 const value = {.x = (float)r / 255.0f, .y = (float)g / 255.0f, .z = (float)b / 255.0f, .w = a};
   return value;
}

static ImVec4 colorDoubleFull(float const a)
{
   return color(0xC4, 0x38, 0x9F, a);
}

static ImVec4 colorDoubleDark(float const a)
{
   return color(0x31, 0x01, 0x38, a);
}

static ImVec4 colorTripleFull(float const a)
{
   return color(0x21, 0xFF, 0x43, a);
}

static ImVec4 colorTripleDark(float const a)
{
   return color(0x06, 0xCC, 0x94, a);
}

static ImVec4 colorTripleLight(float const a)
{
   return color(0x51, 0x99, 0x58, a);
}

static void appSetStyle(ImGuiStyle *const style)
{
   style->Colors[ImGuiCol_Text] = colorTripleFull(1.0f);
   style->Colors[ImGuiCol_TextDisabled] = colorTripleDark(1.0f);

   style->Colors[ImGuiCol_WindowBg] = colorDoubleDark(0.8f);
   style->Colors[ImGuiCol_PopupBg] = colorDoubleDark(0.75f);

   style->Colors[ImGuiCol_TitleBgActive] = colorTripleLight(1.0f);
   style->Colors[ImGuiCol_FrameBg] = colorTripleLight(0.54f);

   style->Colors[ImGuiCol_FrameBgHovered] = colorTripleDark(0.4f);
   style->Colors[ImGuiCol_FrameBgActive] = colorTripleDark(0.67f);
   style->Colors[ImGuiCol_CheckMark] = colorTripleDark(1.0f);
   style->Colors[ImGuiCol_SliderGrabActive] = colorTripleDark(1.0f);
   style->Colors[ImGuiCol_Button] = colorTripleDark(0.4f);
   style->Colors[ImGuiCol_ButtonHovered] = colorTripleDark(1.0f);
   style->Colors[ImGuiCol_Header] = colorTripleLight(0.70f);
   style->Colors[ImGuiCol_HeaderHovered] = colorTripleDark(0.8f);
   style->Colors[ImGuiCol_HeaderActive] = colorTripleDark(1.0f);
   style->Colors[ImGuiCol_ResizeGrip] = colorTripleDark(0.25f);
   style->Colors[ImGuiCol_ResizeGripHovered] = colorTripleDark(0.67f);
   style->Colors[ImGuiCol_ResizeGripActive] = colorTripleDark(0.95f);
   style->Colors[ImGuiCol_TextSelectedBg] = colorTripleDark(0.35f);

   style->Colors[ImGuiCol_Tab] = colorTripleLight(0.54f);
   style->Colors[ImGuiCol_TabHovered] = colorTripleLight(0.75f);
   style->Colors[ImGuiCol_TabSelected] = colorTripleLight(1.0f);

   style->Colors[ImGuiCol_SliderGrab] = colorDoubleFull(1.0f);
   style->Colors[ImGuiCol_ButtonActive] = colorDoubleFull(1.0f);
   style->Colors[ImGuiCol_SeparatorHovered] = colorDoubleFull(0.78f);
   style->Colors[ImGuiCol_SeparatorActive] = colorDoubleFull(1.0f);

   style->WindowRounding = 0.0f;
}

static struct String newString(char const *const base)
{
   struct String result;
   char const *const safeBase = (base != NULL) ? base : "";
   result.len = SDL_strlen(safeBase);
   result.text = SDL_strdup(safeBase);
   result.text[result.len] = 0x00;
   return result;
}

static struct String newStringFallback(char const *const base, char const *const fallback)
{
   return (base != NULL) ? newString(base) : newString(fallback);
}

SDL_AppResult SDL_AppInit(void **const appstate, int const argc, char *argv[])
{
   (void)argc;
   (void)argv;

   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return appFailSDL("failed to initialize SDL");
   }

   SDL_SetAppMetadata("InkyBlackness - HackEd", REPO_SHORT_VERSION, "io.github.inkyblackness.hacked");
   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return appFailSDL("failed to initialize SDL");
   }
   struct HackEdApp *app = SDL_malloc(sizeof(*app));
   if (app == NULL)
   {
      return appFailSDL("failed to allocate application memory");
   }
   SDL_zerop(app);
   {
      SDL_DisplayMode const *const mode = SDL_GetDesktopDisplayMode(SDL_GetPrimaryDisplay());
      app->originalDesktopWidth = mode->w;
      app->originalDesktopHeight = mode->h;

      app->folderHome = newStringFallback(SDL_GetUserFolder(SDL_FOLDER_HOME), "n/a");
      app->folderDocuments = newStringFallback(SDL_GetUserFolder(SDL_FOLDER_DOCUMENTS), "n/a");
      app->folderApp = newStringFallback(SDL_GetBasePath(), "n/a");

      app->folderPreferences.text = SDL_GetPrefPath("InkyBlackness", "HackEd");
      if (app->folderPreferences.text == NULL)
      {
         app->folderPreferences = newString("n/a");
      }
      app->folderPreferences.len = SDL_strlen(app->folderPreferences.text);
   }
   if (!SDL_CreateWindowAndRenderer(
          "InkyBlackness - HackEd - " REPO_LONG_VERSION, 320, 200, SDL_WINDOW_RESIZABLE | SDL_WINDOW_HIGH_PIXEL_DENSITY, &app->window, &app->renderer))
   {
      // TODO free app properly
      SDL_free(app);
      return appFailSDL("failed to create window/renderer");
   }
   *appstate = app;

   if (!SDL_SetRenderVSync(app->renderer, 1))
   {
      SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "SDL_SetRenderVSync failed: %s", SDL_GetError());
   }
   SDL_SetRenderDrawBlendMode(app->renderer, SDL_BLENDMODE_BLEND); // Ensure blend mode is set on all platforms
   SDL_SetHint(SDL_HINT_MAIN_CALLBACK_RATE, "60.0");

   CIMGUI_CHECKVERSION();
   ImGui_CreateContext(NULL);
   ImGuiIO *io = ImGui_GetIO();
   io->ConfigFlags |= ImGuiConfigFlags_NavEnableKeyboard;
   io->ConfigFlags |= ImGuiConfigFlags_DockingEnable;
   io->IniFilename = NULL;

   ImGuiStyle *style = ImGui_GetStyle();
   appSetStyle(style);
   style->FontScaleDpi = appGetBaseUIScale(app->window);

   cImGui_ImplSDL3_InitForSDLRenderer(app->window, app->renderer);
   cImGui_ImplSDLRenderer3_Init(app->renderer);

   return SDL_APP_CONTINUE;
}

SDL_AppResult SDL_AppEvent(void *const appstate, SDL_Event *const event)
{
   struct HackEdApp const *const app = appstate;

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

   SDL_ConvertEventToRenderCoordinates(app->renderer, event);

   cImGui_ImplSDL3_ProcessEvent(event);

   return SDL_APP_CONTINUE;
}

static void appClearBackground(SDL_Renderer *const renderer)
{
   SDL_SetRenderDrawColor(renderer, 0x00, 0x00, 0x00, 0xFF);
   SDL_RenderClear(renderer);
}

SDL_AppResult SDL_AppIterate(void *const appstate)
{
   struct HackEdApp *const app = appstate;
   SDL_AppResult appResult = SDL_APP_CONTINUE;

   // static bool showSystemInfo = false;

   {
      cImGui_ImplSDLRenderer3_NewFrame();
      cImGui_ImplSDL3_NewFrame();
      ImGui_NewFrame();

      if (ImGui_BeginMainMenuBar())
      {
         if (ImGui_BeginMenu("File"))
         {
            if (ImGui_MenuItem("Exit"))
            {
               // TODO: route exit request through system
               appResult = SDL_APP_SUCCESS;
            }
            ImGui_EndMenu();
         }
         if (ImGui_BeginMenu("About"))
         {
            ImGui_Checkbox("Demo Window", &app->showDemoWindow);
            ImGui_EndMenu();
         }
         ImGui_EndMainMenuBar();
      }

      if (app->showDemoWindow)
      {
         ImGui_ShowDemoWindow(&app->showDemoWindow);
      }

      ImGui_Render();
   }

   appClearBackground(app->renderer);

   SDL_SetRenderScale(app->renderer, ImGui_GetIO()->DisplayFramebufferScale.x, ImGui_GetIO()->DisplayFramebufferScale.y);
   cImGui_ImplSDLRenderer3_RenderDrawData(ImGui_GetDrawData(), app->renderer);
   SDL_RenderPresent(app->renderer);

   return appResult;
}

void SDL_AppQuit(void *const appstate, SDL_AppResult const result)
{
   (void)result;

   struct HackEdApp *const app = appstate;

   if (app)
   {
      SDL_Log("Quitting");

      cImGui_ImplSDLRenderer3_Shutdown();
      cImGui_ImplSDL3_Shutdown();
      ImGui_DestroyContext(ImGui_GetCurrentContext());

      SDL_DestroyRenderer(app->renderer);
      SDL_DestroyWindow(app->window);
      SDL_free(app->folderHome.text);
      SDL_free(app->folderDocuments.text);
      SDL_free(app->folderApp.text);
      SDL_free(app->folderPreferences.text);
      SDL_free(app);
   }
}
