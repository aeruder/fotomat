// +build go1.13

package main

import (
	"flag"
	"testing"
)

func init() {
	testing.Init()
	flag.Parse()
	lateInit()
}
