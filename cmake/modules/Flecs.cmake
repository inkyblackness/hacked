FetchContent_Declare(
        flecs
        GIT_REPOSITORY https://github.com/SanderMertens/flecs.git
        # Matching v4.1.6
        GIT_TAG fb55f3c25660425cfe1bc4cf5e6bff8b3f18a9b8
        EXCLUDE_FROM_ALL
)
SET(FLECS_SHARED OFF CACHE BOOL "Disabled shared libraries" FORCE)
FetchContent_MakeAvailable(flecs)

target_compile_definitions(flecs_static
        PUBLIC
        FLECS_CUSTOM_BUILD
        FLECS_APP
        # FLECS_LOG disabled until bug about mismatching log types is fixed about mismatching function signatures
        FLECS_MODULE
        FLECS_NO_OS_API_IMPL # to force os_api_impl to not include anything, due to it getting re-enabled
        __COSMOCC__ # to force os_api.c to assume execinfo is not available
        FLECS_SYSTEM
        FLECS_STATS
        FLECS_TIMER
)
if (DOS)
    target_compile_definitions(flecs_static PUBLIC ECS_TARGET_FREEBSD) # to force os_api.h to include stdlib.h
endif ()
