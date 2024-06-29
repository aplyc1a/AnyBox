#!/bin/bash
BASE_DIR=`pwd`
BUILD_DIR=${BASE_DIR}/output

YARA_NAME="yara"
YARA_SRC="/shared/yara"
YARA_INCLUDE="${YARA_SRC}/libyara/include"
YARA_LIBARIES="${YARA_SRC}/.libs"

mkdir -p ${BUILD_DIR}

function build_linux {
    echo "Building for linux..."
    for dir in ${BASE_DIR}/cmd/*; do
        if [ -d "$dir" ]; then
            # 获取子项目名称（目录名称）
            project_name=$(basename $dir)

            # 输出当前正在构建的子项目
            echo "Building ${project_name}..."

            if [ "$dir" = "${BASE_DIR}/cmd/pudu" ]; then
                (cd ${YARA_SRC}&&make clean;./configure --enable-static --disable-shared;make)
                (cd "$dir" && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-I${YARA_INCLUDE}" CGO_LDFLAGS="-L${YARA_LIBARIES} -l${YARA_NAME} -lm  -lcrypto -lssl -lz -lzstd" go build -a -ldflags '-s -w --extldflags "-static  -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}")
                continue
            fi

            # 进入子项目目录
            (cd "$dir" && GOOS=linux GOARCH=amd64 go build -ldflags '-s -w --extldflags "-static -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}")
        fi
    done
}

function build_linux_arm64 {
    echo "Building for linux arm64..."
    for dir in ${BASE_DIR}/cmd/*; do
        if [ -d "$dir" ]; then
            # 获取子项目名称（目录名称）
            project_name=$(basename $dir)

            # 输出当前正在构建的子项目
            echo "Building ${project_name}..."

            if [ "$dir" = "${BASE_DIR}/cmd/pudu" ]; then
                (cd ${YARA_SRC}&&make clean;./configure --host=aarch64-linux-gnu --enable-static --disable-shared;make)
                (cd "$dir" && GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CGO_CFLAGS="-I${YARA_INCLUDE}" CGO_LDFLAGS="-L${YARA_LIBARIES} -l${YARA_NAME} -lm  -lcrypto -lssl -lz -lzstd" go build -a -ldflags '-s -w --extldflags "-static  -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}.arm64")
                continue
            fi

            # 进入子项目目录
            (cd "$dir" && GOOS=linux GOARCH=arm64 go build -ldflags '-s -w --extldflags "-static -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}.arm64")
        fi
    done
}

function build_linux_arm7 {
    echo "Building for linux arm7..."
    for dir in ${BASE_DIR}/cmd/*; do
        if [ -d "$dir" ]; then
            # 获取子项目名称（目录名称）
            project_name=$(basename $dir)

            # 输出当前正在构建的子项目
            echo "Building ${project_name}..."

            if [ "$dir" = "${BASE_DIR}/cmd/pudu" ]; then
                (cd ${YARA_SRC}&&make clean;./configure --host=arm-linux-gnueabihf --enable-static --disable-shared;make)
                (cd "$dir" && GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CGO_CFLAGS="-I${YARA_INCLUDE}" CGO_LDFLAGS="-L${YARA_LIBARIES} -l${YARA_NAME} -lm  -lcrypto -lssl -lz -lzstd" go build -a -ldflags '-s -w --extldflags "-static  -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}.arm7")
                continue
            fi

            # 进入子项目目录
            (cd "$dir" && GOOS=linux GOARCH=arm GOARM=7 go build -ldflags '-s -w --extldflags "-static -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}.arm7")
        fi
    done
}

function build_windows {
    echo "Building for windows..."
    for dir in ${BASE_DIR}/cmd/*; do
        if [ -d "$dir" ]; then
            # 获取子项目名称（目录名称）
            project_name=$(basename $dir)

            # 输出当前正在构建的子项目
            echo "Building ${project_name}..."

            if [ "$dir" = "${BASE_DIR}/cmd/pudu" ]; then
                (cd ${YARA_SRC}&&make clean;./configure --host=x86_64-w64-mingw32 --enable-static --disable-shared;make)
                (cd "$dir" && CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-I${YARA_INCLUDE}" CGO_LDFLAGS="-L${YARA_LIBARIES} -lyara" go build -ldflags '-s -w --extldflags "-static -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}.exe")
                continue
            fi

            # 进入子项目目录
            (cd "$dir" && CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags '-s -w --extldflags "-static -fpic"' -buildmode=pie -o "${BUILD_DIR}/${project_name}.exe")
        fi
    done
}

function clean {
    echo "Cleaning..."
    rm -rf "$BUILD_DIR"
}

function main {
    if [ ! -d "$BUILD_DIR" ]; then
        mkdir "$BUILD_DIR"
    fi

    case $1 in
        "linux")
            build_linux
            ;;
        "linux_arm64")
            build_linux_arm64
            ;;
        "linux_arm7")
            build_linux_arm7
            ;;
        "windows")
            build_windows
            ;;
        "all")
            build_linux
            build_linux_arm7
            build_linux_arm64
            build_windows
            ;;
        "clean")
            clean
            ;;
        *)
            echo "Usage: ./build.sh {linux|linux_arm64|linux_arm7|windows|clean}"
            exit 1
            ;;
    esac
    #cp ${BASE_DIR}/web ${BUILD_DIR}/ -r
}

main $1
