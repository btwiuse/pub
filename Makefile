all:
	go mod tidy
	go generate
	./run.ts

fmt:
	go fmt
	deno fmt
