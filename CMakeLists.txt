cmake_minimum_required(VERSION 3.28)
project(Hacked LANGUAGES C)

add_executable(hacked
        main.c
)

install(TARGETS hacked
        RUNTIME DESTINATION .
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
