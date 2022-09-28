# Set up variables
DIST_DIR := dist
APP_NAME := $(shell go list -m)
APP_VERSION := $(file < VERSION)
APP_COMMIT := $(shell git rev-list -1 HEAD)
LD_FLAGS := "-X 'mtsaver/app.BuildVersion=${APP_VERSION}' -X 'mtsaver/app.BuildCommit=${APP_COMMIT}'"

fn_GO_BUILD = GOOS=$(1) GOARCH=$(2) go build -o ${DIST_DIR}/$(3) -ldflags=${LD_FLAGS} main.go ;\
7z a ${DIST_DIR}/${APP_NAME}-${APP_VERSION}-$(4).7z -mx9 ./${DIST_DIR}/$(3) ;\
rm ${DIST_DIR}/$(3)


all: clean build

.PHONY: build
build: build-windows build-linux


.PHONY: build-windows
build-windows: ${DIST_DIR}
	$(call fn_GO_BUILD,windows,amd64,mtsaver.exe,win64)
	$(call fn_GO_BUILD,windows,386,mtsaver.exe,win32)


.PHONY: build-linux
build-linux: ${DIST_DIR}
	$(call fn_GO_BUILD,linux,amd64,mtsaver.exe,linux64)
	$(call fn_GO_BUILD,linux,386,mtsaver.exe,linux32)


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
