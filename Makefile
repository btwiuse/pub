build-linux:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o staticlib/linux/x86_64/libpub.a -buildmode=c-archive .

all:
	go mod tidy
	go generate
	./run.ts

zigbuild:
	cargo zigbuild --release --target x86_64-unknown-linux-musl
	ldd ./target/x86_64-unknown-linux-musl/release/main
	ldd ./target/x86_64-unknown-linux-musl/release/static

rust:
	cargo build --release
	ldd ./target/release/main
	ldd ./target/release/static
	# env LD_LIBRARY_PATH=$PWD ./target/release/main
	# ./target/release/static

fmt:
	go fmt
	deno fmt

clean:
	rm -rf staticlib*
	rm libpub*
