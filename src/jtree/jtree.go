package jtree

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// each node of tree contains:
//   size = total size of:
//             - files
//             - subdirectories (not their contents)
//   dirs = map of:
//          - key = subdir name (relative)
//          - value = pointer to subdir tree node
//
// each node of tree DOES NOT contain:
//   it's own name
//   it's own directory size
//
// size totals include:
//     regular files
//     sym links (not what they point to)

type Tree struct {
	size int64
	dirs map[string]*Tree
}

func Descend(dir, path string) *Tree {
	//log.Println("  dir =", dir)
	log.Println("descend =", path)  // path only used for debugging

	tree := &Tree{0, make(map[string]*Tree)}

	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s, %s\n", path, err)
		return tree
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			//log.Println("file is dir =", file.Name())

			subtree := Descend(file.Name(), filepath.Join(path, file.Name()))
			tree.dirs[file.Name()] = subtree
			tree.size += subtree.size
		}
		tree.size += file.Size()
	}

	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}

	return tree
}

/*
func oldDescend(dir string) *Tree {
	// top node has to be treated differently:
	// 1. to get the size of the top directory itself (because the
	//    recursive function only includes size of child directories)
	// 2. to get the real path to the top directory

	if err := os.Chdir(dir); err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("real start dir = ", dir)

	fileinfo, err := os.Lstat(".")
	if err != nil {
		log.Fatal(err)
	}

	// $tree->{'size'} = -s $startDir;
	//$tree->{'dirs'}->{$startDir} = descend($startDir, $startDir);

	tree := &Tree{fileinfo.Size(), make(map[string]*Tree)}
	tree.dirs[dir] = Descend(dir, dir)

	return tree
}
*/
func _DumpTree(tree *Tree, depth uint) {

	for dir, subtree := range tree.dirs {
		fmt.Printf("%d,%d,%s\n", subtree.size, depth, dir)
		_DumpTree(subtree, depth+1)
	}
}
func DumpTree(tree *Tree) {
	fmt.Println("#size,depth,name")
	fmt.Printf("%d,%d,%s\n", tree.size, 0, "{TOP}")
	_DumpTree(tree, 1)
}
