#!/bin/sh
# With apologies to the Go Authors(tm). Tabs have done me wrong so many times
# in the past that I just can't open my heart and trust them again. Sorry.
gofmt -w -tabs=false -tabwidth=4 $*
