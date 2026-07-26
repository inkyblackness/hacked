FetchContent_Declare(
        SDL
        GIT_REPOSITORY https://github.com/libsdl-org/SDL
        # This is a semi-random tag after 3.14.12, but before a release
        # It includes MS DOS support.
        GIT_TAG 855cbec702f246661ff00a0bce9e0683012840c2
        EXCLUDE_FROM_ALL
)

# Ensure only static libray is used
set(SDL_SHARED OFF CACHE BOOL "Disabled SDL shared lib" FORCE)
set(SDL_STATIC ON CACHE BOOL "Enabled SDL static lib" FORCE)

# Disable unused build functions
set(SDL_TEST OFF CACHE BOOL "Disable SDL Test" FORCE)
set(SDL_INSTALL OFF CACHE BOOL "Disable SDL installation" FORCE)

# Disable unused features
set(SDL_DISKAUDIO OFF CACHE BOOL "Disable SDL support for disk audio" FORCE)
set(SDL_CAMERA OFF CACHE BOOL "Disable SDL support for camera" FORCE)
set(SDL_GPU_OPENXR OFF CACHE BOOL "Disable SDL OpenXR support" FORCE)
set(SDL_X11_XSCRNSAVER OFF CACHE BOOL "Disable SDL xscreensaver dependency" FORCE)
set(SDL_X11_XTEST OFF CACHE BOOL "Disable SDL xtest dependency" FORCE)

FetchContent_MakeAvailable(SDL)
