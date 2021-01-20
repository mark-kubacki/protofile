// Package protofile implements temporary files that don't appear in
// filesystem namespace, and get discarded, unless "hard-linked" before
// closing them.
//
// Only works on Linux and file systems supporting O_TMPFILE.
package protofile
