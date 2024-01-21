package main

//go:generate env CGO_ENABLED=1 go build -v -o libpub.so -buildmode=c-shared .
//go:generate env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o staticlib/linux/x86_64/libpub.a -buildmode=c-archive .
//go:generate env CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -v -o staticlib/linux/aarch64/libpub.a -buildmode=c-archive .
//go:generate tar -cz staticlib -f staticlib.tgz
