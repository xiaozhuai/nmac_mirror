#!/usr/bin/env bash

cd frontend
yarn
yarn build
cd ../backend
bindata ./public/...

echo "Build nmac_mirror_darwin_amd64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -i -o ./build/nmac_mirror_darwin_amd64 .

echo "Build nmac_mirror_linux_amd64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -o ./build/nmac_mirror_linux_amd64 .

echo "Build nmac_mirror_linux_arm ..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -i -o ./build/nmac_mirror_linux_arm .

echo "Build nmac_mirror_linux_arm64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -i -o ./build/nmac_mirror_linux_arm64 .

echo "Build nmac_mirror_win_386.exe ..."
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -i -o ./build/nmac_mirror_win_386.exe .

echo "Build nmac_mirror_win_amd64.exe ..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -i -o ./build/nmac_mirror_win_amd64.exe .

echo "Build nmac_mirror_freebsd_amd64 ..."
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -i -o ./build/nmac_mirror_freebsd_amd64 .

rm -rf zip
mkdir zip
cp -f ./build/nmac_mirror_darwin_amd64 ./zip/
cp -f ./build/nmac_mirror_linux_amd64 ./zip/
cp -f ./build/nmac_mirror_linux_arm ./zip/
cp -f ./build/nmac_mirror_linux_arm64 ./zip/
cp -f ./build/nmac_mirror_win_386.exe ./zip/
cp -f ./build/nmac_mirror_win_amd64.exe ./zip/
cp -f ./build/nmac_mirror_freebsd_amd64 ./zip/
cp -f config.yaml ./zip/

cd zip
zip nmac_mirror_darwin_amd64.zip config.yaml nmac_mirror_darwin_amd64
zip nmac_mirror_linux_amd64.zip config.yaml nmac_mirror_linux_amd64
zip nmac_mirror_linux_arm.zip config.yaml nmac_mirror_linux_arm
zip nmac_mirror_linux_arm64.zip config.yaml nmac_mirror_linux_arm64
zip nmac_mirror_win_386.zip config.yaml nmac_mirror_win_386.exe
zip nmac_mirror_win_amd64.zip config.yaml nmac_mirror_win_amd64.exe
zip nmac_mirror_freebsd_amd64.zip config.yaml nmac_mirror_freebsd_amd64
