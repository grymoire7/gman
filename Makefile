# Makefile for gman project
# TODO: Add source dependencies
#       Add code formatting

ifeq ($(mode),debug)
	GOBUILD = go build -gcflags "-N -l"
	GOTEST = go test -c
else
	GOBUILD = go build
	GOTEST = go test
endif

# Go parameters
GOCMD=go
GOPATH=$(PWD)
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

BUILD=gman
TEST=test_terminal test_man2md

.PHONY:all
all: clean $(BUILD) $(TEST)

test: $(TEST)

gman:
	GOPATH=$(GOPATH) && $(GOBUILD) gman

test_terminal:
	cd $(GOPATH)/src/blackfriday && $(GOTEST) -run Term

test_man2md:
	cd $(GOPATH)/src/man2md && $(GOTEST) -run Man

clean:
	-rm -f gman



