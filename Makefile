# Set up variables
APP_NAME := mtsaver
DIST_DIR := dist

ifeq (${OS},Windows_NT)
	BUILD_TIME := $(shell date)
else
	BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
endif

MODULE_NAME := $(shell go list)
APP_VERSION := $(file < VERSION)
APP_COMMIT := $(shell git rev-list -1 HEAD)
LD_FLAGS := "-X '${MODULE_NAME}/app.BuildVersion=${APP_VERSION}' -X '${MODULE_NAME}/app.BuildCommit=${APP_COMMIT}' -X '${MODULE_NAME}/app.BuildTime=${BUILD_TIME}'"

fn_GO_BUILD = GOOS=$(1) GOARCH=$(2) go build -o ${DIST_DIR}/$(3) -ldflags=${LD_FLAGS} main.go ;\
7z a ${DIST_DIR}/${APP_NAME}-${APP_VERSION}-$(4).7z -mx9 ./${DIST_DIR}/$(3)


all: clean build-dist

.PHONY: build-dist
build-dist: build-linux64 build-linux32 build-windows32 build-windows64
	rm -f ${DIST_DIR}/${APP_NAME}
	rm -f ${DIST_DIR}/${APP_NAME}.exe


.PHONY: build-windows32
build-windows32: clean ${DIST_DIR}
	$(call fn_GO_BUILD,windows,386,${APP_NAME}.exe,win32)

.PHONY: build-windows64
build-windows64: clean ${DIST_DIR}
	$(call fn_GO_BUILD,windows,amd64,${APP_NAME}.exe,win64)

.PHONY: build-linux32
build-linux32: clean ${DIST_DIR}
	$(call fn_GO_BUILD,linux,386,${APP_NAME},linux32)

.PHONY: build-linux64
build-linux64: clean ${DIST_DIR}
	$(call fn_GO_BUILD,linux,amd64,${APP_NAME},linux64)


${DIST_DIR}:
	mkdir ${DIST_DIR}


.PHONY: version
version:
	@echo Version from 'VERSION' file: ${APP_VERSION}


.PHONY: clean
clean:
	rm -rf ${DIST_DIR}
