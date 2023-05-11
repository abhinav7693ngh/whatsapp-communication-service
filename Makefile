BINARY_NAME=multiBot
.DEFAULT_GOAL := run

build:
	go build -o ./bin/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf bin

dep: vet
	go mod download

vet:
	go vet

run:
	make clean
	make dep
	make build

fast:
	make clean
	make dep
	make build
	./bin/${BINARY_NAME}
