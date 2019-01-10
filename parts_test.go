// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package parts_test

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"

	"github.com/apatters/go-parts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Pseudo-constants
var (
	testDataPaths = []string{
		"testdata/test.conf",
		"testdata/etc",
		"testdata/usr/lib",
	}
	testDataAllFiles = []string{
		"testdata/etc/10-both.conf",
		"testdata/etc/10-only-etc.conf",
		"testdata/etc/20-only-etc.conf",
		"testdata/etc/30-symlink.conf",
		"testdata/test.conf",
		"testdata/usr/lib/10-executable.sh",
		"testdata/usr/lib/10-only-lib.conf",
		"testdata/usr/lib/20-only-lib.conf",
		"testdata/usr/lib/40-noconf",
		"testdata/usr/lib/adir",
		"testdata/usr/lib/nodigits.conf",
	}
	testDataDefaultFiles = []string{
		"testdata/etc/10-both.conf",
		"testdata/etc/10-only-etc.conf",
		"testdata/etc/20-only-etc.conf",
		"testdata/etc/30-symlink.conf",
		"testdata/test.conf",
		"testdata/usr/lib/10-executable.sh",
		"testdata/usr/lib/10-only-lib.conf",
		"testdata/usr/lib/20-only-lib.conf",
		"testdata/usr/lib/40-noconf",
		"testdata/usr/lib/nodigits.conf",
	}
	testDataExecutables = []string{
		"testdata/usr/lib/10-executable.sh",
	}
	testDataDirs = []string{
		"testdata/usr/lib/adir",
	}
	testDataConfigFiles = []string{
		"testdata/etc/10-both.conf",
		"testdata/etc/10-only-etc.conf",
		"testdata/etc/20-only-etc.conf",
		"testdata/etc/30-symlink.conf",
		"testdata/test.conf",
		"testdata/usr/lib/10-only-lib.conf",
		"testdata/usr/lib/20-only-lib.conf",
		"testdata/usr/lib/nodigits.conf",
	}
	testDataSymlinks = []string{
		"testdata/etc/30-symlink.conf",
	}
)

func TestWalkAllFiles(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.ModeType,
		parts.ModePerm,
		parts.DefaultRegExpFilter)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Logf("testDataPaths: %s", testDataPaths)
	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("err: %v", err)
	t.Logf("fileNames: %s", fileNames)
	require.NoError(t, err)

	assert.EqualValuesf(
		t,
		testDataAllFiles,
		fileNames,
		"Expected file list does not match result.")
}

func TestWalkAllFilesWithN(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.ModeType,
		parts.ModePerm,
		parts.DefaultRegExpFilter)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Logf("testDataPaths: %s", testDataPaths)
	p := parts.NewParts(testDataPaths, config)
	n := 5
	fileNames, err := p.Readdirnames(n)
	t.Logf("n = %d", n)
	t.Logf("err: %v", err)
	t.Logf("fileNames: %s", fileNames)
	require.NoError(t, err)

	assert.EqualValuesf(
		t,
		testDataAllFiles[0:5],
		fileNames,
		"Expected file list does not match result.")
}

func TestWalkAllFilesReversed(t *testing.T) {
	config, err := parts.NewConfig(
		true,
		parts.ModeType,
		parts.ModePerm,
		parts.DefaultRegExpFilter)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Logf("testDataPaths: %s", testDataPaths)
	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("err: %v", err)
	t.Logf("fileNames: %s", fileNames)
	require.NoError(t, err)

	reversedTestDataAllFiles := make([]string, len(testDataAllFiles))
	copy(reversedTestDataAllFiles, testDataAllFiles)
	sort.Sort(sort.Reverse(sort.StringSlice(reversedTestDataAllFiles)))

	assert.EqualValuesf(
		t,
		reversedTestDataAllFiles,
		fileNames,
		"Expected file list does not match result.")
}

func TestWalkDefaults(t *testing.T) {
	p := parts.NewParts(testDataPaths, nil)
	fileNames, err := p.Readdirnames(0)
	t.Logf("config: %v", p.Config)
	t.Logf("err: %v", err)
	t.Logf("testDataPaths: %s", testDataPaths)
	t.Logf("fileNames: %s", fileNames)

	require.NoError(t, err)
	assert.EqualValuesf(
		t,
		testDataDefaultFiles,
		fileNames,
		"Expected file list does not match result.")
}

func TestWalkExecutables(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.ExecutableModeTypeFilter,
		parts.ExecutableModePermFilter,
		parts.DefaultRegExpFilter)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("err: %v", err)
	t.Logf("testDataPaths: %s", testDataPaths)
	t.Logf("fileNames: %s", fileNames)
	require.NoError(t, err)

	assert.EqualValuesf(
		t,
		testDataExecutables,
		fileNames,
		"Expected file list does not match result.")
}

func TestWalkDirectories(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.ModeDir,
		parts.ModePerm,
		parts.DefaultRegExpFilter)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("err: %v", err)
	t.Logf("testDataPaths: %s", testDataPaths)
	t.Logf("fileNames: %s", fileNames)
	require.NoError(t, err)

	assert.EqualValuesf(
		t,
		testDataDirs,
		fileNames,
		"Expected file list does not match result.")
}

func TestWalkConfig(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.DefaultModeTypeFilter,
		parts.DefaultModePermFilter,
		`\.conf$`)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("config: %v", p.Config)
	t.Logf("err: %v", err)
	t.Logf("testDataPaths: %s", testDataPaths)
	t.Logf("fileNames: %s", fileNames)

	require.NoError(t, err)
	assert.EqualValuesf(
		t,
		testDataConfigFiles,
		fileNames,
		"Expected file list does not match result.")
}

func TestReadConfigRaw(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.DefaultModeTypeFilter,
		parts.DefaultModePermFilter,
		`\.conf$`)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("config: %v", p.Config)
	t.Logf("err: %v", err)
	t.Logf("testDataPaths: %s", testDataPaths)
	t.Logf("fileNames: %s", fileNames)

	require.NoError(t, err)
	require.EqualValuesf(
		t,
		testDataConfigFiles,
		fileNames,
		"Expected file list does not match result.")

	defer p.Close()
	expectedContents := ""
	for _, fileName := range testDataConfigFiles {
		expectedContents += filepath.Base(fileName) + "\n"
	}
	contents := ""
	var numBytes int
	buf := make([]byte, 8)
	for {
		numBytes, err = p.Read(buf)
		if numBytes > 0 {
			contents += string(buf[0:numBytes])
		}
		if err != nil {
			break
		}
	}
	t.Logf("err: %v", err)
	require.Equal(t, err, io.EOF, "Returned non-EOF error")
	t.Logf("expectedContents:\n%s", expectedContents)
	t.Logf("contents:\n%s", contents)
	require.EqualValues(t, expectedContents, contents)
}

func TestReadConfigWithReader(t *testing.T) {
	config, err := parts.NewConfig(
		false,
		parts.DefaultModeTypeFilter,
		parts.DefaultModePermFilter,
		`\.conf$`)
	t.Logf("config: %v", config)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	p := parts.NewParts(testDataPaths, config)
	fileNames, err := p.Readdirnames(0)
	t.Logf("config: %v", p.Config)
	t.Logf("err: %v", err)
	t.Logf("testDataPaths: %s", testDataPaths)
	t.Logf("fileNames: %s", fileNames)

	require.NoError(t, err)
	require.EqualValuesf(
		t,
		testDataConfigFiles,
		fileNames,
		"Expected file list does not match result.")

	expectedContents := ""
	for _, fileName := range testDataConfigFiles {
		expectedContents += filepath.Base(fileName) + "\n"
	}

	defer p.Close()
	b, err := ioutil.ReadAll(p)
	t.Logf("err: %v", err)
	require.NoError(t, err)
	require.NotEmptyf(t, b, "Readall() resulted in empty buffer")
	contents := string(b)
	t.Logf("expectedContents:\n%s", expectedContents)
	t.Logf("contents:\n%s", contents)
	require.EqualValues(t, expectedContents, contents)
}

func TestNoDir(t *testing.T) {
	p := parts.NewParts([]string{"/notexist"}, nil)
	fileNames, err := p.Readdirnames(0)
	t.Logf("filenames: %v", fileNames)
	t.Logf("err: %v", err)
	assert.Error(t, err)
	assert.Empty(t, fileNames)

	defer p.Close()
	b, err := ioutil.ReadAll(p)
	t.Logf("err: %v", err)
	t.Logf("b: %q", string(b))
	assert.Error(t, err)
	assert.Empty(t, b)
}

func TestStatMode(t *testing.T) {
	fileName := testDataExecutables[0]
	mode, err := parts.StatMode(fileName)
	t.Logf("fileName: %s", fileName)
	t.Logf("mode: %s", mode)
	t.Logf("err: %v", err)
	assert.True(t, mode.IsExecutable())
	assert.NoError(t, err)

	fileName = testDataDirs[0]
	mode, err = parts.StatMode(fileName)
	t.Logf("fileName: %s", fileName)
	t.Logf("mode: %s", mode)
	t.Logf("err: %v", err)
	assert.True(t, mode.IsDir())
	assert.NoError(t, err)

	fileName = testDataConfigFiles[0]
	mode, err = parts.StatMode(fileName)
	t.Logf("fileName: %s", fileName)
	t.Logf("mode: %s", mode)
	t.Logf("err: %v", err)
	assert.True(t, mode.IsRegular())
	assert.NoError(t, err)
}

func TestLstatMode(t *testing.T) {
	symlink := testDataSymlinks[0]
	mode, err := parts.LstatMode(symlink)
	t.Logf("symlink: %s", symlink)
	t.Logf("mode: %s", mode)
	t.Logf("err: %v", err)
	assert.True(t, mode&parts.ModeSymlink != 0)
	assert.NoError(t, err)

	mode, err = parts.StatMode(symlink)
	t.Logf("symlink: %s", symlink)
	t.Logf("mode: %s", mode)
	t.Logf("err: %v", err)
	assert.True(t, mode&parts.ModeSymlink == 0)
	assert.NoError(t, err)
}
