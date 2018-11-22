// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

/*
Anticycle is a tool for static code analysis which search for
dependency cycles. It scans recursively all source files and
parses theirs dependencies. Anticycle does not compile the code,
so it is ideal for searching for complex, difficult to debug cycles.

Anticycle consists of three parts:

1. Main package which exists in cmd/anticycle directory. It compiles to CLI binary

2. Public API which is importable for other applications.

3. Internal, private interface which should not be imported and use directly. It may change at any moment without warnings.
*/
package anticycle

// Imports for godoc to follow packages
import (
	_ "github.com/anticycle/anticycle/pkg/anticycle"
	_ "github.com/anticycle/anticycle/pkg/model"
	_ "github.com/anticycle/anticycle/pkg/serialize"
)
