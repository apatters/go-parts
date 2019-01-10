// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package parts

// Partser implements all functions needed for a ReadCloser and adds
// the os.File.Readdirnames method.
type Partser interface {
	Readdirnames(int) ([]string, error)
	Read([]byte) (int, error)
	Close() error
}
