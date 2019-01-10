// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package parts

import (
	"os"
)

// FileMode is an extension of the os.FileMode struct. It adds a
// ModeRegular bit to make testing for a regular file explicit.
type FileMode os.FileMode

const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        = FileMode(os.ModeDir)        // d: is a directory
	ModeAppend     = FileMode(os.ModeAppend)     // a: append-only
	ModeExclusive  = FileMode(os.ModeExclusive)  // l: exclusive use
	ModeTemporary  = FileMode(os.ModeTemporary)  // T: temporary file; Plan 9 only
	ModeSymlink    = FileMode(os.ModeSymlink)    // L: symbolic link
	ModeDevice     = FileMode(os.ModeDevice)     // D: device file
	ModeNamedPipe  = FileMode(os.ModeNamedPipe)  // p: named pipe (FIFO)
	ModeSocket     = FileMode(os.ModeSocket)     // S: Unix domain socket12
	ModeSetuid     = FileMode(os.ModeSetuid)     // u: setuid
	ModeSetgid     = FileMode(os.ModeSetgid)     // g: setgid
	ModeCharDevice = FileMode(os.ModeCharDevice) // c: Unix character device, when ModeDevice is set
	ModeSticky     = FileMode(os.ModeSticky)     // t: sticky
	ModeRegular    = 1 << (31 - iota)            // f: regular file

	// Mask for the type bits.
	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeRegular

	// Unix permission bits
	ModePerm FileMode = 0777
)

func (m FileMode) String() string {
	const str = "dalTLDpSugctf"
	var buf [32]byte // Mode is uint32.
	w := 0
	for i, c := range str {
		if m&(1<<uint(32-1-i)) != 0 {
			buf[w] = byte(c)
			w++
		}
	}
	if w == 0 {
		buf[w] = '-'
		w++
	}
	const rwx = "rwxrwxrwx"
	for i, c := range rwx {
		if m&(1<<uint(9-1-i)) != 0 {
			buf[w] = byte(c)
		} else {
			buf[w] = '-'
		}
		w++
	}
	return string(buf[:w])
}

// IsDir reports whether m describes a directory.
func (m FileMode) IsDir() bool {
	return m&ModeDir != 0
}

// IsRegular reports whether m describes a regular file.
func (m FileMode) IsRegular() bool {
	switch {
	case m&ModeRegular != 0:
		return true
	// Maintain conversion compatibility with os.FileMode.
	case m&ModeType == 0:
		return true
	default:
		return false
	}
}

// Perm returns the permission bits of the file mode.
func (m FileMode) Perm() FileMode {
	return m & ModePerm
}

// IsExecutable reports whether m describes an executablexs file.
func (m FileMode) IsExecutable() bool {
	return m.IsRegular() && (m&0111 != 0)
}

// FileInfo extends the os.FileInfo struct.
type FileInfo struct {
	os.FileInfo
}

// Mode returns the file mode bits of the file with an adjustment made
// for regular files.
func (i FileInfo) Mode() FileMode {
	m := i.FileInfo.Mode()
	if m&os.ModeType == 0 {
		return ModeRegular
	}

	return FileMode(m)
}
