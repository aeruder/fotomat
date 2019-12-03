// +build !go1.13

package main

import (
	"flag"
)

func init() {
	// Initialize flags with default values, enable local serving.
	flag.Parse()
	lateInit()
}
