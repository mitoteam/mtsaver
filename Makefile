all: clean build

DIST_DIR := dist

build: ${DIST_DIR}
	go build -o ${DIST_DIR}/mtsaver.exe main.go

${DIST_DIR}:
	mkdir ${DIST_DIR}

.PHONY: clean
clean:
	rm -rf dist
