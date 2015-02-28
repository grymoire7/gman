# Makefile for gman project
# TODO: Add source dependencies

ifeq ($(mode),debug)
	GOBUILD = go build -gcflags "-N -l"
	GOTEST = go test -c
else
	GOBUILD = go build
	GOTEST = go test
endif

# Go parameters
GOCMD=go
export GOPATH=$(PWD)
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GODEP=$(GOTEST) -i
# With apologies to the Go Authors(tm). Tabs have done me wrong so many times
# in the past that I just can't open my heart and trust them again. Sorry.
GOFMT=gofmt -w -tabs=false -tabwidth=4
GOFMTGO=gofmt -w
GOGET=go get
BUILD=gman
TEST=test_terminal test_man2md

.PHONY: clean get fmt $(BUILD) $(TEST)

# This works not because it's called 'default' but because it's first.
default: gman

gman:
	GOPATH=$(GOPATH) $(GOBUILD) gman

test: $(TEST)

fmt:
	$(GOFMT) ./src/gman && $(GOFMTGO) ./src/man2md

get:
	$(GOGET) github.com/grymoire7/docopt.go; \
	$(GOGET) github.com/grymoire7/blackfriday

test_terminal:
	cd $(GOPATH)/src/github.com/grymoire7/blackfriday && $(GOTEST) -run Term

test_man2md:
	cd $(GOPATH)/src/man2md && $(GOTEST) -run Man

clean:
	-rm -f gman



