package jtree

import (
	"testing"
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"path/filepath"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

const empty_dir_size int64 = 4096

func TestJTree(t *testing.T) {
	// Setup
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	//defer os.RemoveAll(dir)

	fmt.Println(dir)
    os.MkdirAll(filepath.Join(dir, "dir1"), 0755)
	os.MkdirAll(filepath.Join(dir, "dir2"), 0755)
	bytes := []byte("123456789")  // 9 bytes
	test_file_size := int64(len(bytes))
	err = ioutil.WriteFile(filepath.Join(dir, "file1"), bytes, 0644)
	err = ioutil.WriteFile(filepath.Join(dir, "dir2", "file1"), bytes, 0644)
	err = ioutil.WriteFile(filepath.Join(dir, "dir2", "file2"), bytes, 0644)
	tree := Descend(dir)
	log.Println(tree)

	t.Run("CheckStartNode", func(t *testing.T) {
		var expected_size int64 = empty_dir_size
		if tree.size != expected_size {
			t.Fatalf("expected top size to be %v, got %v", expected_size, tree.size)
		}

		if len(tree.dirs) != 1 {
			t.Fatalf("expected 1 dir in start node, got %v", len(tree.dirs))
		}
	})

	// get first (and only) child node from start node
	var topdir *Tree = nil
	for _, node := range tree.dirs {
		topdir = node
		break
	}

	t.Run("CheckTopDir", func(t *testing.T) {
		var expected_size int64 = empty_dir_size * 2 + test_file_size // 8201

		if topdir.size != expected_size {
			t.Fatalf("expected top size to be %v, got %v", expected_size, topdir.size)
		}

		if len(topdir.dirs) != 2 {
			t.Fatalf("expected 2 dirs in top node, got %v", len(tree.dirs))
		}
	})

	t.Run("CheckDir1", func(t *testing.T) {
		var expected_size int64 = 0

		dir1 := topdir.dirs["dir1"]
		if dir1.size != expected_size {
			t.Fatalf("expected dir1 size to be %v, got %v", expected_size, dir1.size)
		}

		if len(dir1.dirs) != 0 {
			t.Fatalf("expected 0 dirs in top node, got %v", len(dir1.dirs))
		}
	})

	t.Run("CheckDir2", func(t *testing.T) {
		var expected_size int64 = test_file_size * 2

		dir2 := topdir.dirs["dir2"]
		if dir2.size != expected_size {
			t.Fatalf("expected dir1 size to be %v, got %v", expected_size, dir2.size)
		}

		if len(dir2.dirs) != 0 {
			t.Fatalf("expected 0 dirs in top node, got %v", len(dir2.dirs))
		}
	})
}

