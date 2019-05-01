.PHONY: clean build test all
GO_SOURCES = $(shell find . -type f -name '*.go')

all: test build

build: artifactory/io/projectriff/command/io.projectriff.command

test:
	go test -v ./...

artifactory/io/projectriff/command/io.projectriff.command: buildpack.toml $(GO_SOURCES)
	rm -fR $@ 							&& \
	./ci/package.sh						&& \
	mkdir $@/latest 					&& \
	tar -C $@/latest -xzf $@/*/*.tgz

clean:
	rm -fR artifactory/
	rm -fR dependency-cache/
