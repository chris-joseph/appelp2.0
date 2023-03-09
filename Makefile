
run: build-go

build-go:
	go build -buildmode=c-shared -o ./vendor/out/main.dll ./vendor/main.go