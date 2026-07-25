cmake_minimum_required(VERSION 3.28)
project(Hacked)

cmake_policy(SET CMP0077 NEW) # Abort if variables are not clearly set

set(CMAKE_C_STANDARD 23)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_CXX_STANDARD 17) # GoogleTest requires at least C++17
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)

FetchContent_Declare(
        googletest
        GIT_REPOSITORY https://github.com/google/googletest.git
        GIT_TAG 58d77fa8070e8cec2dc1ed015d66b454c8d78850 # release-1.12.1
        OVERRIDE_FIND_PACKAGE
)
# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

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
