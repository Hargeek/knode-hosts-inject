run:
	go run main.go

build: clean
	go build -v -o bin/server main.go

buildx:
	go build -ldflags="-s -w" -o /bin/server main.go

.PHONY: clean
clean:
	go clean
