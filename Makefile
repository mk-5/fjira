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

build_windows:
	mkdir -p ${BINARY_DIR}
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}.exe cmd/fjira-cli/main.go
	chmod +x ${BINARY_NAME}.exe

run:
	./${BINARY_NAME}

test:
	go test ./internal/...

test_coverage:
	go test -coverpkg=./... -covermode=count -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func coverage.out

release:
	goreleaser release --skip-publish --snapshot --rm-dist

clean:
	rm -rf ${BINARY_DIR}
	rm -rf dist
