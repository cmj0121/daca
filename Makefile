PROJ=daca
SRC=$(wildcard **/*.go)

SUBDIR=

.PHONY: all gofmt test run clean

all: gofmt test

gofmt:
	gofmt -s -w $(SRC)

test: gofmt
	go test -v ./...

clean:
	rm -rf $(RELEASE)

.PHONY: clean

clean:
