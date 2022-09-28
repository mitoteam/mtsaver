APP_VERSION := $(file < VERSION)
DIST_DIR := dist

all: clean build

.PHONY: build
build: ${DIST_DIR}
	go build -o ${DIST_DIR}/mtsaver.exe -ldflags="-X 'mtsaver/app.BuildVersion=${APP_VERSION}'" main.go


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
