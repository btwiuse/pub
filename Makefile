tidy:
	git config --global --add safe.directory '*'
	go mod tidy

img:
	docker build -t btwiuse/pub:dev .
	docker push btwiuse/pub:dev

devcontainer:
	docker build -t btwiuse/pub:dev .
	docker push btwiuse/pub:devcontainer

build-linux: tidy
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o staticlib/linux/x86_64/libpub.a -buildmode=c-archive .

build-linux-arm64: tidy
	which aarch64-linux-gnu-gcc || sudo pacman -Sy aarch64-linux-gnu-gcc --noconfirm
	env CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -v -o staticlib/linux/aarch64/libpub.a -buildmode=c-archive .

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
