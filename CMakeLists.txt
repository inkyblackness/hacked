cmake_minimum_required(VERSION 3.28)
project(Hacked)

cmake_policy(SET CMP0077 NEW) # Abort if variables are not clearly set

# Setting a C-Standard automatically reduces the available compilers; For example, MSVC is out.
# See here for a good overview: https://en.cppreference.com/c/compiler_support/23
set(CMAKE_C_STANDARD 23)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_CXX_STANDARD 17) # GoogleTest requires at least C++17
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)
include(CTest)

if (NOT DOS)
    # As long as djgpp does not support building GoogleTest, there is no testing capability on MS DOS.
    include(cmake/modules/TestSupport.cmake)
endif ()
include(cmake/modules/SimpleDirectMediaLayer.cmake)
include(cmake/modules/Nuklear.cmake)
include(cmake/modules/BuildRules.cmake)
include(cmake/modules/Version.cmake)

file(GLOB_RECURSE LIB_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/main/source/*.c")
add_library(hacked-core STATIC ${LIB_CORE_SOURCES})
target_include_directories(hacked-core PUBLIC "${PROJECT_SOURCE_DIR}/core/main/include")
target_link_libraries(hacked-core
        PRIVATE
        BuildRules::common
)

if (NOT DOS)
    file(GLOB_RECURSE TEST_SUPPORT_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/test-support/source/*.cpp")
    add_library(hacked-core-test-support STATIC ${TEST_SUPPORT_CORE_SOURCES})
    target_include_directories(hacked-core-test-support PUBLIC "${PROJECT_SOURCE_DIR}/core/test-support/include")
    target_link_libraries(hacked-core-test-support
            PUBLIC
            hacked-core
            GTest::gmock
            PRIVATE
            BuildRules::common
    )

    file(GLOB_RECURSE TEST_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/test/source/*.cpp")
    add_executable(hacked-core-test ${TEST_CORE_SOURCES})
    target_link_libraries(hacked-core-test
            hacked-core
            hacked-core-test-support
            BuildRules::common
            GTest::gmock_main
    )
    gtest_add_tests(TARGET hacked-core-test EXTRA_ARGS --gtest_shuffle SKIP_DEPENDENCY)
endif ()

file(GLOB_RECURSE LIB_NUKLEAR_SDL_BRIDGE_SOURCES "${PROJECT_SOURCE_DIR}/nuklear-sdl-bridge/main/source/*.c")
add_library(hacked-nuklear-sdl-bridge STATIC ${LIB_NUKLEAR_SDL_BRIDGE_SOURCES})
target_include_directories(hacked-nuklear-sdl-bridge PUBLIC "${PROJECT_SOURCE_DIR}/nuklear-sdl-bridge/main/include")
target_link_libraries(hacked-nuklear-sdl-bridge
        PUBLIC
        hacked-core
        Nuklear
        SDL3::SDL3-static
)

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
    set_target_properties(hacked PROPERTIES SUFFIX ".html") # Tell Emscripten to also create .html for this (needed?)
endif ()
target_link_libraries(hacked
        hacked-core
        hacked-nuklear-sdl-bridge
        Version::all
        BuildRules::common
        SDL3::SDL3-static
)

include(cmake/modules/InstallationRules.cmake)
