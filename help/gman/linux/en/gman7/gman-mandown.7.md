Gman-gmd(7) -- Gman manual markdown format
==========================================

## SYNOPSIS

The Gman Markdown format (.gmd) is Markdown with a few help system extensions.
The format is similar to, and was inspired by, the
[Ronn format](https://github.com/rtomayko/ronn).

## DESCRIPTION

One of the features of the gman(1) command is to convert text in a simple
format, a version of Markdown(1), to user consummable formats like UNIX style
man pages or interactively browsed HTML.

Not all groff(1) or Markdown(1) typesetting features can be expressed in the
gmd syntax.

## TITLE
## SECTION HEADINGS
Frequently used section headings:
```
## TLDR
## SYNOPSIS
## DESCRIPTION
## OPTIONS
## SYNTAX
## ENVIRONMENT
## RETURN VALUES
## STANDARDS
## SECURITY CONSIDERATIONS
## BUGS
## HISTORY
## AUTHOR
## COPYRIGHT
## SEE ALSO
```

## INLINE MARKUP
## DEFINITION LISTS
## LINKS
## SEE ALSO
ronn(1), ronn-format(7), markdown(7), groff(7)

