// Copyright 2018 The Anticycle Authors. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

/*
Anticycle is a tool for static code analysis which search for
dependency cycles. It scans recursively all source files and
parses theirs dependencies. Anticycle does not compile the code,
so it is ideal for searching for complex, difficult to debug cycles.

Anticycle consists of three parts.
1.	main package which exists in cmd/anticycle directory
	it compiles to binary which can be used within terminal.
2.	public API which is meant to be import by other applications.
3.	internal, private functions and structures which should not
	be imported and use directly. They may change at any moment.
*/
package anticycle
