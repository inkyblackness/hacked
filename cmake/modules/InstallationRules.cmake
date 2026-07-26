# Installation rules specify what needs to be put into a package for distribution.
# They are relevant for the execution of
# cmake --install <build-dir> --prefix <distribution-dir>

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
