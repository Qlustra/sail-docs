// Command catalog demonstrates a declaration-driven catalog preview and
// explicit reads from a caller-selected live workspace.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := demonstrate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
