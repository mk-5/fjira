BINARY_DIR=out/bin
BINARY_NAME=${BINARY_DIR}/fjira

all: clean install test build

build_run: clean build run

install:
	go mod vendor

build:
	mkdir -p ${BINARY_DIR}
	go build -o ${BINARY_NAME} cmd/fjira-cli/main.go
	chmod +x ${BINARY_NAME}

run:
	./${BINARY_NAME}

test:
	go test ./internal/...

clean:
	rm -rf ${BINARY_DIR}
