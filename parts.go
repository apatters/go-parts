// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

/*
Package parts is used to traverse or read files in one or more
directories using Debian's run-parts layout. Filters can be used to
only include:

    * Files that match a regular expression.
    * Files that match a file type mask, e.g., regular files, directories, etc.
    * Files that match a file permission mask, e.g., executables.

Directories are specified in a slice. Files from directories
appearing earlier in the slice that have the same base name as a
file from a directory appearing later in the slice take
precedence. e.g., using the following directory layouts:

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

The etc/10-both.conf file is used and the /usr/lib/10-both.conf is
ignored if etc is configured ahead of /usr/lib.

After duplicates are resolved, the files are lexically sorted
(optionally in reverse order).

Using the above layout and using a ".*\.conf" file name filter we
can use parts.Readdirnames() to list the files in run-parts order,
e.g.,

    etc/10-both.conf
    etc/10-only-etc.conf
    usr/lib/10-only-lib.conf
    etc/20-only-etc.conf
    usr/lib/20-only-lib.conf
    etc/30-symlink.conf
    usr/lib/nodigits.conf
    test.conf

Or we can read the concatenated contents of all these files using
parts.Read().
*/
package parts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	DefaultModeTypeFilter    = ModeRegular
	DefaultModePermFilter    = ModePerm
	ExecutableModePermFilter = 0111
	ExecutableModeTypeFilter = ModeRegular
	DefaultRegExpFilter      = ".*"
)

// Config is used to pass parameter to NewConfig.
type Config struct {
	Reverse        bool
	ModeTypeFilter FileMode
	ModePermFilter FileMode
	RegExpFilter   *regexp.Regexp
}

// NewConfig constructor. Can fail if regular expressions do not
// compile.
func NewConfig(reverse bool, modeTypeFilter FileMode, modePermFilter FileMode, regExpFilter string) (*Config, error) {
	regExp, err := regexp.Compile(regExpFilter)
	if err != nil {
		return nil, fmt.Errorf("parts: %s", err)
	}
	return &Config{
		Reverse:        reverse,
		ModeTypeFilter: modeTypeFilter,
		ModePermFilter: modePermFilter,
		RegExpFilter:   regExp,
	}, nil
}

// NewDefaultConfig returns a default Config constructor.
func NewDefaultConfig() *Config {
	return &Config{
		Reverse:        false,
		ModeTypeFilter: DefaultModeTypeFilter,
		ModePermFilter: DefaultModePermFilter,
		RegExpFilter:   regexp.MustCompile(DefaultRegExpFilter),
	}
}

// readState tracks the current state of reading the parts
// directory. It is mostly used to close directory files during Read()
// operations.
type readState struct {
	Files  []io.ReadCloser
	Reader io.Reader
}

// Parts encapsulates data and functions used to process "run-parts"
// directories.
type Parts struct {
	Paths     []string
	Config    *Config
	readState *readState
}

// NewParts is the Parts constructor. A default configuration is used
// if config is nil.
func NewParts(paths []string, config *Config) *Parts {
	if config == nil {
		config = NewDefaultConfig()
	}
	return &Parts{
		Paths:     paths,
		Config:    config,
		readState: nil,
	}
}

// Readdirnames returns a list of files in paths that follow the
// "run-parts" naming convention.
func (p *Parts) Readdirnames(n int) ([]string, error) {
	foundNames := make(map[string]string)
	for _, path := range p.Paths {
		mode, err := StatMode(path)
		if err != nil {
			return []string{}, fmt.Errorf("parts: %s", err)
		}
		switch {
		case mode.IsDir():
			dir, err := os.Open(path)
			if err != nil {
				return []string{}, fmt.Errorf("parts: %s", err)
			}
			defer dir.Close()
			fileNames, err := dir.Readdirnames(0)
			if err != nil {
				return []string{}, fmt.Errorf("parts: %s", err)
			}
			for _, fileName := range fileNames {
				fullPath := filepath.Join(path, fileName)
				mode, err = StatMode(fullPath)
				if err != nil {
					return []string{}, fmt.Errorf("parts: %s", err)
				}
				if _, ok := foundNames[fileName]; ok {
					continue
				}
				if p.filter(fileName, mode, p.Config.RegExpFilter) {
					foundNames[fileName] = fullPath
				}
			}
		default:
			if _, ok := foundNames[filepath.Base(path)]; ok {
				continue
			}
			if p.filter(filepath.Base(path), mode, nil) {
				foundNames[filepath.Base(path)] = path
			}
		}
	}
	names := make([]string, 0, len(foundNames))
	for _, val := range foundNames {
		names = append(names, val)
	}
	if p.Config.Reverse {
		sort.Sort(sort.Reverse(pathsByBasename(names)))
	} else {
		sort.Sort(pathsByBasename(names))
	}

	switch {
	case n == 0:
		return names, nil
	case n < len(names):
		return names[0:n], nil
	default:
		return names, nil
	}
}

// Read reads the contents of the parts directory into buffer b.
func (p *Parts) Read(b []byte) (int, error) {
	if p.readState == nil {
		// Initialize
		foundFiles, err := p.Readdirnames(0)
		if err != nil {
			return 0, err
		}
		// Create a reader for each file and stuff it away.
		p.readState = new(readState)
		p.readState.Files = make([]io.ReadCloser, 0, len(foundFiles))
		for _, fileName := range foundFiles {
			file, err := os.Open(fileName)
			if err != nil {
				return 0, err
			}
			p.readState.Files = append(p.readState.Files, file)
		}
		readers := make([]io.Reader, 0, len(p.readState.Files))
		for _, reader := range p.readState.Files {
			readers = append(readers, reader)
		}
		p.readState.Reader = io.MultiReader(readers...)
	}

	bytesRead, err := p.readState.Reader.Read(b)

	return bytesRead, err
}

// Close closes all files opened by Read.
func (p *Parts) Close() error {
	var err error
	if p.readState == nil {
		return nil
	}
	for _, reader := range p.readState.Files {
		if reader == nil {
			continue
		}
		tmpErr := reader.Close()
		if tmpErr != nil && err != nil {
			err = tmpErr
		}
	}
	p.readState = nil

	return err
}

// filter returns true if the file name and mode matches the filtering
// criteria (name regexp, perms, and mode).
func (p *Parts) filter(name string, mode FileMode, regExp *regexp.Regexp) bool {
	if regExp != nil && !regExp.MatchString(name) {
		return false
	}
	if mode&p.Config.ModePermFilter == 0 {
		return false
	}
	if mode&p.Config.ModeTypeFilter == 0 {
		return false
	}

	return true
}

// StatMode returns the FileMode for the named path. If there is an error,
// it will be of type *PathError.
func StatMode(name string) (FileMode, error) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return 0, err
	}

	var mode FileMode
	if fileInfo.Mode().IsRegular() {
		mode = FileMode(fileInfo.Mode()) | ModeRegular
	} else {
		mode = FileMode(fileInfo.Mode())
	}

	return mode, nil
}

// LstatMode returns the FileMode for the named path. If the path is a
// symbolic link, the returned FileMode describes the symbolic
// link. If there is an error, it will be of type *PathError.
func LstatMode(name string) (FileMode, error) {
	fileInfo, err := os.Lstat(name)
	if err != nil {
		return 0, err
	}

	var mode FileMode
	if fileInfo.Mode().IsRegular() {
		mode = FileMode(fileInfo.Mode()) | ModeRegular
	} else {
		mode = FileMode(fileInfo.Mode())
	}

	return mode, nil
}

type pathsByBasename []string

func (paths pathsByBasename) Len() int {
	return len(paths)
}

func (paths pathsByBasename) Swap(i, j int) {
	paths[i], paths[j] = paths[j], paths[i]
}

func (paths pathsByBasename) Less(i, j int) bool {
	return strings.Compare(filepath.Base(paths[i]), filepath.Base(paths[j])) < 0
}
