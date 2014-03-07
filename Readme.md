<img src="gmanhat.png" align="right"/>
# gman(1) - A better help system

## Summary
Gman is a help system with some modern options.
Sources are in the easy-to-edit Markdown format.
Among other things, this gets you inline images in
the web browsable version.
Here are a few features by example:

### Just show the Summary section:
    alias tldr='gman -s Summary'
    tldr tar

### Just show particular options:
    gman "rsync -a -v -m"

### Start http server for interactive browsing:
    gman -browse

### Query existing man pages:
    gman -q man rsync

## Gman roadmap

### Status
Working on version 0.1.

### Version 0.1
Minimal features for something that works.

* Terminal renderer for Markdown (70%)
* Help page section extraction (100%)
* Compressed help page support (100%)
* Help page contribution guidelines (0%)
* Options extraction (50%)

### Version 0.2
* Inline Ronn-style extensions for parser
* Allow multiple languages and OSes
* Install script
* Bug fixes, help pages, etc.

### Version 0.3
* Allow use of existing man pages.
* Provide tool to convert man pages to gman pages.
* HTTP server option for interactive browsing.
* Examples with images
* Bug fixes, help pages, etc.

### Version 0.4-0.9
* Determine criteria for Version 1.0
* Bug fixes, help pages, etc.


## References
* Markdown processing library [Blackfriday](https://github.com/russross/blackfriday).
* Markdown extensions as described by [Ronn format](https://github.com/rtomayko/ronn).
* Initial project inspiration from [TLDR](https://github.com/rprieto/tldr).


