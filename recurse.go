package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

// RecursiveCompute is the main function, it recurses into directories
// and builds up the tree of files
func RecursiveCompute(c chan bool, node *TreeNode, searchPath string) (err error) {
	node.Path = searchPath

	c <- true
	listing, err := ioutil.ReadDir(searchPath)
	<-c
	if err != nil {
		fmt.Println(searchPath, err)
		return
	}

	var wg sync.WaitGroup
	for _, info := range listing {
		wg.Add(1)
		go func(info os.FileInfo) {
			if info.IsDir() {
				dir := NewNode(info.Name())
				dir.Parent = node
				sub := path.Join(searchPath, info.Name())
				RecursiveCompute(c, &dir, sub)
				node.Insert(info.Name(), &dir)
				node.AddSize(dir.GetSize())
			} else {
				fileEntry := NewFileEntry(info.Name(), uint64(info.Size()))
				node.Insert(info.Name(), &fileEntry)
				node.AddSize(fileEntry.GetSize())
			}
			wg.Done()
		}(info)
	}

	wg.Wait()
	return
}
