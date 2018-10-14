// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar"
	"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/bar/baz"
	"github.com/anticycle/anticycle/internal/pkg/scan/testdata/nocycle/pkg/foo"
)

func main() {
	fmt.Printf("foo %v bar %v baz %v", foo.Foo{}, bar.Bar{}, baz.Baz{})
}
