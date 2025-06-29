ROOT := $(shell pwd)

VERSION=$(shell git describe --tags --long)
TARGET=$(ROOT)/main/gobase
all:build

.PHONY:dep
dep:
	go mod tidy

.PHONY:build
build:dep
	cd main && CGO_ENABLED=1 CGO_CFLAGS="-DSQLITE_ENABLE_RTREE -DSQLITE_THREADSAFE=1" go build -ldflags "-X github.com/open-cmi/gobase/internal/commands.Version=${VERSION} -s -w" -o $(TARGET) main.go

BUILDDIR?=/usr/local
.PHONY:install
install:
	mkdir -p ${BUILDDIR}/bin
	cp -rfp ${TARGET} ${BUILDDIR}/bin/

.PHONY:clean
clean:
	rm -r build/*
