cmake_minimum_required(VERSION 3.28)
project(Hacked LANGUAGES C)

add_executable(hacked
        main.c
)

install(TARGETS hacked
        RUNTIME DESTINATION .
)
