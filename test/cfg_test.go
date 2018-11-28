// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package test

import (
	"flag"
)

// Update flag for sanity and acceptance tests
var update = flag.Bool("update", false, "update .golden files")
