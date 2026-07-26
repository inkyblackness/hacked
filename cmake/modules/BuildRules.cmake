# Starting off with a fresh project, enabling all the warnings is a lofty goal, yet will help in the long run.
add_library(build_rules_warnings_all INTERFACE)
add_library(BuildRules::warnings::all ALIAS build_rules_warnings_all)
target_compile_options(build_rules_warnings_all INTERFACE
        $<$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>,$<CXX_COMPILER_ID:GNU>>:-Wall -Wextra -Wpedantic -Werror>
        # MSVC: /WX for treating warnings as errors breaks for various 3rd-party includes, thus not included
        $<$<CXX_COMPILER_ID:MSVC>:/W4 /Wall>
)

# Optimizations for speed are tricky, as they would require some sort of performance analysis to prove their worth.
# However, we can assume a certain baseline to use.
add_library(build_rules_optimizations_speed INTERFACE)
add_library(BuildRules::optimizations::speed ALIAS build_rules_optimizations_speed)
target_compile_options(build_rules_optimizations_speed INTERFACE
        $<$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>,$<CXX_COMPILER_ID:GNU>>:-O3>
)
target_link_options(build_rules_optimizations_speed INTERFACE
        # LTO := Link Time Optimization
        $<$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>,$<CXX_COMPILER_ID:GNU>>:-flto>
)

add_library(build_rules_common INTERFACE)
add_library(BuildRules::common ALIAS build_rules_common)
target_link_libraries(build_rules_common
        INTERFACE
        BuildRules::warnings::all
        $<$<CONFIG:Release,RelWithDebInfo>:BuildRules::optimizations::speed>
)
