package gdu

import (
	"io/ioutil"
	"os"
	"path"
	"sync"
)

// RecursiveCompute is the main function, it recurses into directories
// and builds up the tree of files
func RecursiveCompute(node *TreeNode, searchPath string) (err error) {
	listing, err := ioutil.ReadDir(searchPath)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, file := range listing {
		wg.Add(1)

		go func(info os.FileInfo) {
			if info.IsDir() {
				dir := NewNode(info.Name() + "/")
				dir.Parent = node
				go RecursiveCompute(&dir, path.Join(searchPath, info.Name()))
				node.Insert(info.Name()+"/", &dir)
				node.AddSize(dir.GetSize())
			} else {
				fileEntry := NewFileEntry(info.Name(), uint64(info.Size()))
				node.Insert(info.Name(), &fileEntry)
				node.AddSize(fileEntry.GetSize())
			}
			wg.Done()
		}(file)
	}

	wg.Wait()
	return
}
