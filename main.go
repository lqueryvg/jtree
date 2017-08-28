package main

import (
	"flag"
	"fmt"
	"os"
	"jtree"
)


func main() {

	startPath := flag.String("d", "", "directory")
	flag.Parse()
	if len(*startPath) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: -d not specified")
		flag.Usage()
		return
	}

	tree := jtree.Descend(*startPath)

	fmt.Println("#size,depth,name")

	jtree.DumpTree(tree, 1)
}
