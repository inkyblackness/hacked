#!/bin/sh

pushd $(dirname "${0}") > /dev/null
HACKED_BASE=$(pwd -L)
FOLDER_NAME=inkyblackness-hacked

echo "Cleaning output directories..."
rm -rf _build

mkdir -p $HACKED_BASE/_build/linux/$FOLDER_NAME
mkdir -p $HACKED_BASE/_build/win/$FOLDER_NAME


echo "Determining version"

VERSION=$(git describe exact-match --abbrev=0)
if [ $? -ne 0 ]; then
   echo "Not a tagged build, defaulting to revision for version"
   REV=$(git rev-parse --short HEAD)
   VERSION="rev$REV"
fi
echo "Determined version: $VERSION"

echo "Building executables..."
go build -ldflags "-X main.version=$VERSION" -o $HACKED_BASE/_build/linux/$FOLDER_NAME/hacked .
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -ldflags "-X main.version=$VERSION" -ldflags -H=windowsgui -o $HACKED_BASE/_build/win/$FOLDER_NAME/hacked.exe .


echo "Copying resources..."

for os in "linux" "win"
do
   packageDir=$HACKED_BASE/_build/$os/$FOLDER_NAME

   cp $HACKED_BASE/LICENSE $packageDir
   cp -R $HACKED_BASE/_resources/* $packageDir
done

MINGW_BASE=/usr/x86_64-w64-mingw32/bin
for lib in "libgcc_s_seh-1.dll" "libstdc++-6.dll" "libwinpthread-1.dll"
do
   cp $MINGW_BASE/$lib $HACKED_BASE/_build/win/$FOLDER_NAME
done


echo "Creating packages..."

cd $HACKED_BASE/_build/linux
tar -cvzf $HACKED_BASE/_build/$FOLDER_NAME-$VERSION.linux64.tgz ./$FOLDER_NAME

cd $HACKED_BASE/_build/win
zip -r $HACKED_BASE/_build/$FOLDER_NAME-$VERSION.win64.zip .

popd > /dev/null
