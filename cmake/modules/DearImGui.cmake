SET(DearImGuiTag "b48d1afbe8ee8b238e2961dc363a949dd7304e23")
SET(DearImGuiVersion "v1.92.9b-docking")
SET(DCImGuiVersion "v0.21")
SET(DCImGuiHash "df0ea4fdb73eca417c0ee8df9d5383362451c95ea3d3d7536c258820023488f9")

FetchContent_Declare(
        dearimgui
        GIT_REPOSITORY https://github.com/ocornut/imgui
        GIT_TAG ${DearImGuiTag}
        OVERRIDE_FIND_PACKAGE
        EXCLUDE_FROM_ALL
)
FetchContent_MakeAvailable(dearimgui)

FetchContent_Declare(
        dcimgui
        URL https://github.com/dearimgui/dear_bindings/releases/download/DearBindings_${DCImGuiVersion}_ImGui_${DearImGuiVersion}/DearBindings_${DCImGuiVersion}_ImGui_${DearImGuiVersion}.zip
        URL_HASH SHA256=${DCImGuiHash}
        OVERRIDE_FIND_PACKAGE
        EXCLUDE_FROM_ALL
)
FetchContent_MakeAvailable(dcimgui)

add_library(DearImGui STATIC
        ${dcimgui_SOURCE_DIR}/backends/dcimgui_impl_sdl3.cpp
        ${dcimgui_SOURCE_DIR}/backends/dcimgui_impl_sdlrenderer3.cpp
        ${dcimgui_SOURCE_DIR}/dcimgui.cpp
        ${dcimgui_SOURCE_DIR}/dcimgui_internal.cpp
        ${dearimgui_SOURCE_DIR}/backends/imgui_impl_sdl3.cpp
        ${dearimgui_SOURCE_DIR}/backends/imgui_impl_sdlrenderer3.cpp
        ${dearimgui_SOURCE_DIR}/imgui.cpp
        ${dearimgui_SOURCE_DIR}/imgui_demo.cpp
        ${dearimgui_SOURCE_DIR}/imgui_draw.cpp
        ${dearimgui_SOURCE_DIR}/imgui_tables.cpp
        ${dearimgui_SOURCE_DIR}/imgui_widgets.cpp
)
target_include_directories(DearImGui
        PUBLIC
        ${dcimgui_SOURCE_DIR}
        ${dcimgui_SOURCE_DIR}/backends
        ${dearimgui_SOURCE_DIR} # TODO Public as long as there is no local IMGUI_USER_CONFIG
        PRIVATE
        ${dearimgui_SOURCE_DIR}/backends
)
target_compile_definitions(DearImGui
        PUBLIC
        IMGUI_DISABLE_OBSOLETE_FUNCTIONS
)
target_link_libraries(DearImGui PUBLIC SDL3::SDL3-static)
