// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protofile implements temporary files that don't appear in
// filesystem namespace, and get discarded, unless "hard-linked" before
// closing them.
//
// Only works on Linux and file systems supporting O_TMPFILE.
package protofile
