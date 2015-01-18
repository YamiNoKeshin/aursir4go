

// Copyright 2013 - by Jim Lawless
// License: MIT / X11
// See: http://www.mailsend-online.com/license2013.php
//
// Bear with me ... I'm a Go noob.

package main

import (
"flag"
"fmt"
)

// Define a type named "intslice" as a slice of ints
type connections []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (i *connections) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *connections) Set(value string) error {
		*i = append(*i, value )
	return nil
}

var myints connections

func main() {
	flag.Var(&myints, "i", "List of integers")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	} else {
		fmt.Println("Here are the values in 'myints'")
		for i := 0; i < len(myints); i++ {
			fmt.Printf("%s\n", myints[i])
		}
	}
}

