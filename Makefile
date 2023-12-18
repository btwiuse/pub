all:
	go generate
	./run.ts

fmt:
	go fmt
	deno fmt
