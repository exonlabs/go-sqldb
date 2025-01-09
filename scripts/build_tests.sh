#!/bin/bash
cd $(dirname $(readlink -f $0))/..

GO=go

BUILD_PATH=build/tests
if [ ! -z "${GO_BIN}" ] ;then
    GO=${GO_BIN}
    BUILD_PATH=${BUILD_PATH}/${GO_BIN}
fi
BUILD_LINUX_PATH=${BUILD_PATH}/linux
BUILD_WIN_PATH=${BUILD_PATH}/win

# clean build dirs
rm -rf ${BUILD_PATH}
mkdir -m 775 -p ${BUILD_LINUX_PATH} ${BUILD_WIN_PATH}

# linux build
for n in sqldb ;do
    ${GO} test ./pkg/${n} -c -o ${BUILD_LINUX_PATH}/${n}.test
done

# windows build
for n in sqldb ;do
    GOOS=windows GOARCH=amd64 ${GO} test \
        ./pkg/${n} -c -o ${BUILD_WIN_PATH}/${n}.test_64.exe
    GOOS=windows GOARCH=386 ${GO} test \
        ./pkg/${n} -c -o ${BUILD_WIN_PATH}/${n}.test_32.exe
done
