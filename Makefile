.PHONY: clean build test acceptance all
GO_SOURCES = $(shell find . -type f -name '*.go')

PACK=go run github.com/buildpacks/pack/cmd/pack

all: test build acceptance

build: artifactory/io/projectriff/command/io.projectriff.command

test:
	go test -v ./...

acceptance:
	$(PACK) create-builder -b acceptance/testdata/builder.toml projectriff/builder
	docker pull cloudfoundry/build:base-cnb
	docker pull cloudfoundry/run:base-cnb
	GO111MODULE=on go test -v -tags=acceptance ./acceptance

artifactory/io/projectriff/command/io.projectriff.command: buildpack.toml $(GO_SOURCES)
	rm -fR $@ 							&& \
	./ci/package.sh						&& \
	mkdir $@/latest 					&& \
	tar -C $@/latest -xzf $@/*/*.tgz

clean:
	rm -fR artifactory/
	rm -fR dependency-cache/
