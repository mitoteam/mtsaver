APP_VERSION := $(file < VERSION)
APP_COMMIT := $(shell git rev-list -1 HEAD)
DIST_DIR := dist

all: clean build

.PHONY: build
build: ${DIST_DIR}
	go build -o ${DIST_DIR}/mtsaver.exe -ldflags="-X 'mtsaver/app.BuildVersion=${APP_VERSION}' -X 'mtsaver/app.BuildCommit=${APP_COMMIT}'" main.go


${DIST_DIR}:
	mkdir ${DIST_DIR}


.PHONY: run
run:
	${DIST_DIR}/mtsaver.exe info


.PHONY: version
version:
	@echo Version from 'VERSION' file: ${APP_VERSION}


.PHONY: clean
clean:
	rm -rf dist
