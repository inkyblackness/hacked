cmake_minimum_required(VERSION 3.28)
project(Hacked)

cmake_policy(SET CMP0077 NEW) # Abort if variables are not clearly set

set(CMAKE_C_STANDARD 23)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_CXX_STANDARD 17) # GoogleTest requires at least C++17
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)
include(CTest)

if (NOT DOS)
    FetchContent_Declare(
            googletest
            URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
            OVERRIDE_FIND_PACKAGE
            EXCLUDE_FROM_ALL
    )
    # For Windows: Prevent overriding the parent project's compiler/linker settings
    set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
    FetchContent_MakeAvailable(googletest)
    include(GoogleTest)
endif ()

FetchContent_Declare(
        SDL
        GIT_REPOSITORY https://github.com/libsdl-org/SDL
        # This is a semi-random tag after 3.14.12, but before a release
        # It includes MS DOS support.
        GIT_TAG 855cbec702f246661ff00a0bce9e0683012840c2
        EXCLUDE_FROM_ALL
)
set(SDL_SHARED OFF CACHE BOOL "Disabled SDL shared lib" FORCE)
set(SDL_STATIC ON CACHE BOOL "Enabled SDL static lib" FORCE)
set(SDL_TEST OFF CACHE BOOL "Disable SDL Test" FORCE)
set(SDL_INSTALL OFF CACHE BOOL "Disable SDL installation" FORCE)
set(SDL_DISKAUDIO OFF CACHE BOOL "Disable SDL support for disk audio" FORCE)
set(SDL_CAMERA OFF CACHE BOOL "Disable SDL support for camera" FORCE)
set(SDL_X11_XSCRNSAVER OFF CACHE BOOL "Disable SDL xscreensaver dependency" FORCE)
set(SDL_X11_XTEST OFF CACHE BOOL "Disable SDL xtest dependency" FORCE)
set(SDL_GPU_OPENXR OFF CACHE BOOL "Disable SDL OpenXR support" FORCE)
FetchContent_MakeAvailable(SDL)

file(GLOB_RECURSE LIB_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/main/source/*.c")
add_library(hacked-core STATIC ${LIB_CORE_SOURCES})
target_include_directories(hacked-core PUBLIC "${PROJECT_SOURCE_DIR}/core/main/include")

if (NOT DOS)
    file(GLOB_RECURSE TEST_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/test/source/*.cpp")
    add_executable(hacked-core-test ${TEST_CORE_SOURCES})
    target_link_libraries(hacked-core-test
            hacked-core
            GTest::gmock_main
    )
    gtest_add_tests(TARGET hacked-core-test EXTRA_ARGS --gtest_shuffle)
endif ()

if (WIN32)
    enable_language(RC)
    set(APP_GUI_WIN_RESOURCES app/mswin/hacked.rc)
endif ()
set(APP_GUI_SOURCES app/source/main.c ${APP_GUI_WIN_RESOURCES})
if (CMAKE_BUILD_TYPE STREQUAL "Debug")
    add_executable(hacked ${APP_GUI_SOURCES})
else ()
    add_executable(hacked WIN32 ${APP_GUI_SOURCES})
endif ()
if (EMSCRIPTEN)
    set_target_properties(hacked PROPERTIES SUFFIX ".html")
    #set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -s ASSERTIONS=1 -s WASM=1 -s ASYNCIFY -s GL_ENABLE_GET_PROC_ADDRESS=1")
endif ()
target_link_libraries(hacked
        hacked-core
        SDL3::SDL3-static
)

# Toolchains may define extra libraries that should be included in the distribution.
# This is primarily for MinGW, which needs some runtime libraries along with it.
if (DEFINED TOOLCHAIN_DLL_DEPENDENCIES)
    foreach (DLL ${TOOLCHAIN_DLL_DEPENDENCIES})
        execute_process(
                COMMAND ${CMAKE_C_COMPILER} -print-file-name=${DLL}
                OUTPUT_VARIABLE DLL_PATH
                OUTPUT_STRIP_TRAILING_WHITESPACE
        )
        if (EXISTS "${DLL_PATH}" AND NOT "${DLL_PATH}" STREQUAL "${DLL}")
            install(FILES "${DLL_PATH}" DESTINATION .)
        else ()
            message(WARNING "DLL '${DLL}' not resolved - it will be missing in the package")
        endif ()
    endforeach ()
endif ()

install(TARGETS hacked
        RUNTIME DESTINATION .
)

if (EMSCRIPTEN)
    install(FILES
            "$<TARGET_FILE_DIR:hacked>/hacked.js"
            "$<TARGET_FILE_DIR:hacked>/hacked.wasm"
            DESTINATION .
    )
endif ()
