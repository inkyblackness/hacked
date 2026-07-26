find_package(Git)

# This module determines version information based on Git.
# The result are two variables: REPO_SHORT_VERSION and REPO_LONG_VERSION.
# These variables are cached, so they can be extracted also from CI.
# Furthermore, they are exposed as preprocessor defines with library Version::all.

execute_process(COMMAND ${GIT_EXECUTABLE} rev-parse --short HEAD OUTPUT_VARIABLE GIT_REVISION ERROR_QUIET)
if ("${GIT_REVISION}" STREQUAL "")
    # Set information to be unknown, in case the source was taken from an exported archive.
    set(REPO_SHORT_VERSION "unknown" CACHE STRING "")
else ()
    execute_process(COMMAND ${GIT_EXECUTABLE} describe --exact-match --tags OUTPUT_VARIABLE GIT_TAG ERROR_QUIET)
    string(STRIP "${GIT_REVISION}" GIT_REVISION)
    string(STRIP "${GIT_TAG}" GIT_TAG)
    if ("${GIT_TAG}" STREQUAL "")
        set(REPO_SHORT_VERSION "rev${GIT_REVISION}" CACHE STRING "")
    else ()
        set(REPO_SHORT_VERSION "${GIT_TAG}" CACHE STRING "")
    endif ()
endif ()

set(REPO_LONG_VERSION "${REPO_SHORT_VERSION}@${CMAKE_SYSTEM_NAME}-${CMAKE_SYSTEM_PROCESSOR}" CACHE STRING "")
message("Determined repository version: '${REPO_LONG_VERSION}'")

add_library(hacked_version INTERFACE)
add_library(Version::all ALIAS hacked_version)
target_compile_definitions(hacked_version
        INTERFACE
        REPO_SHORT_VERSION="${REPO_SHORT_VERSION}"
        REPO_LONG_VERSION="${REPO_LONG_VERSION}"
)
