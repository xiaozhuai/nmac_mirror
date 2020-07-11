#!/usr/bin/env bash

cd frontend
yarn
yarn build
cd ../backend
bindata ./public/...
go build -i -o ./build/nmac_mirror .
