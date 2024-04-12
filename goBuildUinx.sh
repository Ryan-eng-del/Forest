#!/bin/bash

function build_linux() {
    echo =================================
    echo ==========Build Linux ======
    echo =================================
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64
    echo now the CGO_ENABLED:
    go env CGO_ENABLED
    echo now the GOOS:
    go env GOOS
    echo now the GOARCH:
    go env GOARCH
    go build -o "${build_path}/${bin_file}" main.go
}

function build_mac() {
    echo =================================
    echo ==========Build Mac ======
    echo =================================
    export CGO_ENABLED=0
    export GOOS=darwin
    export GOARCH=amd64
    echo now the CGO_ENABLED:
    go env CGO_ENABLED
    echo now the GOOS:
    go env GOOS
    echo now the GOARCH:
    go env GOARCH
    go build -o "${build_path}/${bin_file}" main.go
}

function build_windows() {
    echo =================================
    echo ==========Build Windows ======
    echo =================================
    export CGO_ENABLED=1
    export GOOS=windows
    export GOARCH=amd64
    echo now the CGO_ENABLED:
    go env CGO_ENABLED
    echo now the GOOS:
    go env GOOS
    echo now the GOARCH:
    go env GOARCH
    bin_file="${bin_file}.exe"
    go build -o "${build_path}/${bin_file}" main.go
}


function build() {
  build_os

  if [ -d "${workspace}/build" ];then
    rm -rf "${workspace}/build"
  fi


  build_workspace="${workspace}/build"
  build_dir_name="forest-gateway-${version}-"

  for os in ${arr_build_os[*]}
  do
    build_path="${build_workspace}/${build_dir_name}${os}"
    echo $build_path
    echo "create build dir ${build_path}"
    mkdir -p "${build_path}/install"
    cp -r conf "${build_path}"
    echo ========== Build Forest Gateway ==========
    bin_file="foreset-gateway"
    build_${os}


    echo "pack ${build_dir_name}${os}"
    cd ${build_workspace}
    tar -zcf "${build_dir_name}${os}.tar.gz" ${build_dir_name}${os}
    cd -
  done
}

function build_os() {
  arr_build_os=("windows" "mac" "linux")
}



function main() {
  echo -n  "Input forest gateway release version(1.0.0):"
  workspace=$(cd $(dirname $0); pwd)
  read -a version
  # -E 正则 -o 只输出匹配的字符串
  version=`echo $version |grep -Eo '^[0-9]+.[0-9]+.[0-9]+$'`
  # -n 检测字符串的长度是否不为空
  if [ ! -n "$version" ];then
     echo "input version error example: (1.0.1)"
     exit -1
  fi
  echo $version
  set -e

  # set go proxy
  export GO111MODULE=on && export GOSUMDB=off
  export GOPROXY=https://goproxy.cn

  build
}


main
