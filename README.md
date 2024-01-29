# Linemasker

This is a simple tool that can replace lines with patterns from different files.
(It is just a toy but I need this to simplify a few shell scripts.)
Probably you can achieve the same with some sed magic, but oh well.

## Examples:

`a.txt`

```
burek
1
mesni
1
1
1
```

`b.txt`

```
sirni1
sirni2
```

`linemasker a.txt b.txt`
```
burek
sirni2
mesni
sirni2
sirni1
sirni2
```

`linemasker a.txt b.txt --no-cycle`
```
burek
sirni2
mesni
1
1
1
```

## Installation

Just do:
```
go build ./...
```
