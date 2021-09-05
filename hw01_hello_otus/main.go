package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	const hello = "Hello, OTUS!"
	fmt.Print(stringutil.Reverse(hello))
}
