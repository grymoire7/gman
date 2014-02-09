<img src="gman.1.png" align="right"/>
# gman(1) - A better help system

## Summary
Gman is a help system with some modern options.
Sources are in the easy-to-edit Markdown format.
Here are a few features by example:

**Note:** This help page is currently just an example
and does not (yet) a represent the capabilities of gman.

### Just show the Summary section:
    alias tldr='gman -s Summary'
    tldr tar

### Just show particular options:
    gman "rsync -a -v -z"

### Start http server for interactive browsing:
    gman --browse
    gman -b

### Query existing man pages only:
    gman --query man ipconfig
    gman -q man ipconfig

## Synopsis
gman [-s *section*]
     [-b | --browse]
     [-p | --port *http_port*]
     [-q | --query man]
     [-k | --apropos *regex*]
     *page*

## Options
#### -b, --browse
Start an http server for interactive browsing and launch the default
browser if possible.

#### -k *regex*, --apropos *regex*
Start an http server for interactive browsing and launch the default
browser if possible.

#### -s *section_title*, --section *section_title*
Show only the specified help section. For example, '-s Summary' will display
only the Summary section.

## Gman roadmap
### Status
Working on version 0.1.

### Version 0.1
Minimal features for something that works.

* Blackfriday parser
* Terminal renderer for Blackfriday
* Design docs and roadmap

### Version 0.2
* Section and option extraction
* Inline Ronn-style extensions for Blackfriday
* Allow for multiple languages and OSes.

### Version 0.3
* HTTP server option for interactive browsing.
* Bug fixes, docs, and more gman pages.

### Version 0.4-0.9
* Add support for compressed gman pages.
* Bug fixes, docs, and more gman pages.
* Determine criteria for Version 1.0

### Version ?.?
* Allow use of existing man pages.
* Provide tool to convert man pages to gman pages.

## References
* Markdown processing library [Blackfriday](https://github.com/russross/blackfriday).
* Markdown extensions as described by [Ronn format](https://github.com/rtomayko/ronn).
* Initial project inspiration from [TLDR](https://github.com/rprieto/tldr).


