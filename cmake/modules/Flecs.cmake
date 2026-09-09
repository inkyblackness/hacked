FetchContent_Declare(
        flecs
        GIT_REPOSITORY https://github.com/SanderMertens/flecs.git
        # Matching some time after v4.1.6
        GIT_TAG 16504b2eaf51d5fd8c75b7f2315a49170861643c
        EXCLUDE_FROM_ALL
)
SET(FLECS_SHARED OFF CACHE BOOL "Disabled shared libraries" FORCE)
FetchContent_MakeAvailable(flecs)

target_compile_definitions(flecs_static
        PUBLIC
        FLECS_CUSTOM_BUILD
        FLECS_APP
        FLECS_LOG
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
