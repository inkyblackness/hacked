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
    # As long as djgpp does not support building GoogleTest, there is no testing capability on MS DOS.
    include(cmake/modules/TestSupport.cmake)
endif ()
include(cmake/modules/SimpleDirectMediaLayer.cmake)

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
    set_target_properties(hacked PROPERTIES SUFFIX ".html") # Tell Emscripten to also create .html for this (needed?)
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
    # As per CMake "bug", the additional files have to be manually taken into the installation directory as well.
    # Reference:
    # https://stackoverflow.com/questions/61865545/how-to-install-both-js-and-wasm-files-with-the-cmake-target-for-emscripten#comment140265783_70702138
    install(FILES
            "$<TARGET_FILE_DIR:hacked>/hacked.js"
            "$<TARGET_FILE_DIR:hacked>/hacked.wasm"
            DESTINATION .
    )
endif ()
