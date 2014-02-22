# Makefile for gman project
# TODO: Add source dependencies
#       Add code formatting

# Go parameters
GOCMD=go
GOPATH=$(PWD)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

BUILD=gman
TEST=test_terminal test_man2md

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


