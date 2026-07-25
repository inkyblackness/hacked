cmake_minimum_required(VERSION 3.28)
project(Hacked)

cmake_policy(SET CMP0077 NEW) # Abort if variables are not clearly set

set(CMAKE_C_STANDARD 23)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_CXX_STANDARD 17) # GoogleTest requires at least C++17
#set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)

FetchContent_Declare(
        googletest
        URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
        OVERRIDE_FIND_PACKAGE
)
# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

if (DOS)
    set_target_properties(gtest PROPERTIES
            CXX_STANDARD_REQUIRED OFF
            CXX_EXTENSIONS ON
    )
    set_target_properties(gmock PROPERTIES
            CXX_STANDARD_REQUIRED OFF
            CXX_EXTENSIONS ON
    )
endif ()
# Disable POSIX regular expressions for GoogleTest targets
target_compile_definitions(gtest PUBLIC GTEST_HAS_POSIX_RE=0 GTEST_HAS_STD_WSTRING=0)
target_compile_definitions(gmock PUBLIC GTEST_HAS_POSIX_RE=0 GTEST_HAS_STD_WSTRING=0)

file(GLOB_RECURSE LIB_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/main/source/*.c")
add_library(hacked-core STATIC ${LIB_CORE_SOURCES})
target_include_directories(hacked-core PUBLIC "${PROJECT_SOURCE_DIR}/core/main/include")

file(GLOB_RECURSE TEST_CORE_SOURCES "${PROJECT_SOURCE_DIR}/core/test/source/*.cpp")
add_executable(hacked-core-test ${TEST_CORE_SOURCES})
target_link_libraries(hacked-core-test
        hacked-core
        GTest::gmock_main
)

add_executable(hacked main.c)
target_link_libraries(hacked
        hacked-core
)

install(TARGETS hacked
        RUNTIME DESTINATION .
)
