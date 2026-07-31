FetchContent_Declare(
        nuklear
        GIT_REPOSITORY https://github.com/Immediate-Mode-UI/Nuklear
        # Using an arbitrary tag after v4.13.3, to add newer definitions for SDL3 bridge.
        # TODO: upgrade to a released version.
        GIT_TAG 0dbc52f86404f9e1f26ce0df3015ed23ff54a726
        EXCLUDE_FROM_ALL
)
FetchContent_MakeAvailable(nuklear)

add_library(Nuklear INTERFACE)
target_include_directories(Nuklear INTERFACE ${nuklear_SOURCE_DIR})
