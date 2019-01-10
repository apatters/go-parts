// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package parts_test

import (
	"fmt"
	"io/ioutil"

	"github.com/apatters/go-parts"
)

func ExampleParts_Readdirnames() {
	fmt.Println("List executable files in run-parts order:")
	config, _ := parts.NewConfig(
		false,
		parts.ExecutableModeTypeFilter,
		parts.ExecutableModePermFilter,
		parts.DefaultRegExpFilter)
	p := parts.NewParts(testDataPaths, config)
	names, _ := p.Readdirnames(0)
	for _, name := range names {
		fmt.Printf("\t%s\n", name)
	}

	fmt.Println("\nList files with extension '.conf' in run-parts order:")
	config, _ = parts.NewConfig(
		false,
		parts.DefaultModeTypeFilter,
		parts.DefaultModePermFilter,
		`\.conf$`)
	p = parts.NewParts(testDataPaths, config)
	names, _ = p.Readdirnames(0)
	for _, name := range names {
		fmt.Printf("\t%s\n", name)
	}

	fmt.Println("\nList files with extension '.conf' in reverse run-parts order:")
	config, _ = parts.NewConfig(
		true,
		parts.DefaultModeTypeFilter,
		parts.DefaultModePermFilter,
		`\.conf$`)
	p = parts.NewParts(testDataPaths, config)
	names, _ = p.Readdirnames(0)
	for _, name := range names {
		fmt.Printf("\t%s\n", name)
	}

	// Output:
	// List executable files in run-parts order:
	// 	testdata/usr/lib/10-executable.sh
	//
	// List files with extension '.conf' in run-parts order:
	// 	testdata/etc/10-both.conf
	// 	testdata/etc/10-only-etc.conf
	// 	testdata/usr/lib/10-only-lib.conf
	// 	testdata/etc/20-only-etc.conf
	// 	testdata/usr/lib/20-only-lib.conf
	// 	testdata/etc/30-symlink.conf
	// 	testdata/usr/lib/nodigits.conf
	// 	testdata/test.conf
	//
	// List files with extension '.conf' in reverse run-parts order:
	// 	testdata/test.conf
	// 	testdata/usr/lib/nodigits.conf
	// 	testdata/etc/30-symlink.conf
	// 	testdata/usr/lib/20-only-lib.conf
	// 	testdata/etc/20-only-etc.conf
	// 	testdata/usr/lib/10-only-lib.conf
	// 	testdata/etc/10-only-etc.conf
	// 	testdata/etc/10-both.conf
}

func ExampleParts_Read() {
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
	// Output:
	// 10-both.conf
	// 10-only-etc.conf
	// 10-only-lib.conf
	// 20-only-etc.conf
	// 20-only-lib.conf
	// 30-symlink.conf
	// nodigits.conf
	// test.conf
}

func ExampleStatMode() {
	for _, path := range testDataExecutables {
		mode, _ := parts.StatMode(path)
		fmt.Println(mode)
	}
	// Output: frwxr-xr-x
}

func ExampleLstatMode() {
	for _, path := range testDataSymlinks {
		mode, _ := parts.LstatMode(path)
		fmt.Println(mode)
	}
	// Output: Lrwxrwxrwx
}
