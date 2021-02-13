package main

import (
	"fmt"
	"io"
	"strings"
)

type Prefixer struct {
	writer io.Writer
	prefix string
	first  bool
}

// New creates a new instance of Prefixer.
func NewPrefixWriter(w io.Writer, prefix string) *Prefixer {
	return &Prefixer{
		writer: w,
		prefix: prefix,
	}
}

func (me *Prefixer) Write(p []byte) (n int, err error) {
	// if me.first == false {
	// 	fmt.Fprintf(me.writer, "\n")
	// 	me.first = true
	// }
	parts := strings.Split(string(p), "\n")
	for _, part := range parts {
		fmt.Fprintf(me.writer, "%s %s\n", me.prefix, part)
	}
	return len(p), nil
}
