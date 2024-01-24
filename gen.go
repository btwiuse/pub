package pub

//go:generate env CGO_ENABLED=1 go build -v -o libpub.so -buildmode=c-shared ./libpub
