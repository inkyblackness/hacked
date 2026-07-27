# Building instructions

## What do you want to achieve?

When you want to build the binaries yourself, you typically fall into one of the following categories:

* a) You want to develop/build on your local machine, to run it on that machine.
* b) You want to cross-compile on your local machine, to run it on another architecture.
* c) You want to build all the supported formats.

The necessary tasks of the respective categories are explained below. First, some general information.

## General information

### Supported environments

The project is set up to officially support the following environments:

* Linux (32bit, 64bit)
* MS Windows (32bit, 64bit)
* MacOS (Apple Silicon)
* Internet Browser (WASM)
* MS DOS

> Only what builds by the automated build jobs of this project is supported.
> If there is a permutation of another system that isn't ensured by a build job, it may not work.
>
> The following describes the ideal working setup.

### Dependencies

* [CMake](https://cmake.org/) 3.x, with a minimum version according to `CMakeLists.txt`
  > 4.x should also work, yet may not if compatibilities are removed.
* GNU compiler suite or compatible (Original, MinGW, djgpp; also Clang)
* Internet connection for downloading library dependencies

## How to build?

### You want to develop/build on your local machine and run it on that machine.

#### Simple case

```
# 1) Download library dependencies and prepare build scripts (once, unless project structure is changed):
cmake -B build -S .
# 2) Build binary:
cmake --build build
# 3) Find the executables in the created build/ directory.
```

### You want to cross-compile on your local machine for a different target.

Cross-compilation is considered to be a build that targets a different architecture than what is default for your
machine. This also includes building a 32-bit variant on your 64-bit machine!

For such cases, use and provide
a [CMake toolchain file](https://cmake.org/cmake/help/latest/manual/cmake-toolchains.7.html) for the first command:

```
cmake -B build -S . -DCMAKE_TOOLCHAIN_FILE=<path-to-your-file>
```

Refer to `cmake/toolchains` subdirectory for the ones this project uses. For WASM builds, refer to the build pipeline
for setting up [Emscripten](https://emscripten.org).

### You want to build all the supported formats.

This is essentially the work of the build pipeline. See `.github/workflows/main-build.yml` that combines all the
information above.

On a Linux system, or an MS Windows one with WSL2, you can cover most of the targets.
