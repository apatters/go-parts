# go-parts

Golang Debian run-parts package

[![GoDoc](https://godoc.org/github.com/apatters/go-parts?status.svg)](https://godoc.org/github.com/apatters/go-parts)

Parts is a Go (golang) package that is used to traverse or read files
in one or more directories using Debian's run-parts layout. Filters
can be used to only include:

* Files that match a regular expression.
* Files that match a file type mask, e.g., regular files, directories, etc.
* Files that match a file permission mask, e.g., executables.

Directories are specified in a slice. Files from directories appearing
earlier in the slice that have the same basename as a file from a
directory appearing later in the slice take precedence. Example: using
the following directory layouts:

```
── etc
│   ├── 10-both.conf
│   ├── 10-only-etc.conf
│   ├── 20-only-etc.conf
│   └── 30-symlink.conf -> ../usr/lib/30-symlink.conf
├── test.conf
└── usr
    └── lib
        ├── 10-both.conf
        ├── 10-executable.sh
        ├── 10-only-lib.conf
        ├── 20-only-lib.conf
        ├── 30-symlink.conf
        ├── 40-noconf
        └── nodigits.conf
```

The `etc/10-both.conf` file is used and the `/usr/lib/10-both.conf` is
ignored.

After duplicates are resolved, the files are lexically sorted
(optionally in reverse order).

Using the above layout and using a ".*\.conf" file name filter we
can use parts.Readdirnames() to list the files in run-parts order,
e.g.,

```
etc/10-both.conf
etc/10-only-etc.conf
usr/lib/10-only-lib.conf
etc/20-only-etc.conf
usr/lib/20-only-lib.conf
etc/30-symlink.conf
usr/lib/nodigits.conf
test.conf
```

Or we can read the concatinated contents of all these files using
parts.Read().

Documentation
-------------

Documentation can be found at [GoDoc](https://godoc.org/github.com/apatters/go-parts)

Installation
------------

```bash
$ go get -u github.com/apatters/go-parts
```

## Examples

### Parts.Readdirnames()

``` go
	config, _ := parts.NewConfig(
		false,
		parts.ExecutableModeTypeFilter,
		parts.ExecutableModePermFilter,
		parts.DefaultRegExpFilter)

	p := parts.NewParts(testDataPaths, config)
	names, _ := p.Readdirnames(0)
	for _, name := range names {
		fmt.Println(name)
	}
```

### Parts.Read()

``` go
	config, _ := parts.NewConfig(
		false,
		parts.DefaultModeTypeFilter,
		parts.DefaultModePermFilter,
		`\.conf$`)

	p := parts.NewParts(testDataPaths, config)
	defer p.Close()
	b, _ := ioutil.ReadAll(p)
	contents := string(b)
	fmt.Print(contents)
```

License
-------

MIT.


Thanks
------

Thanks to [Secure64](https://secure64.com/company/) for
contributing this code.
