package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
)

var err error

var root TreeNode
var current *TreeNode

var poundMap = make(map[int]string)

func getPounds(i int) string {
	x, ok := poundMap[i]
	if !ok {
		if i > 0 {
			x = strings.Repeat("#", i)
		} else {
			x = strings.Repeat(" ", -i)
		}
		poundMap[i] = x
	}
	return x
}

func load(ui *UI, node *TreeNode) {
	ui.Table.Clear()
	ui.Table.Path = node.Path
	entries := node.GetEntries()

	if len(entries) > 0 {
		sort.Sort(entries)

		// now entries[0] is guaranteed to be the largest
		largest := entries[0].GetSize()

		if current != &root {
			ui.Table.AddRow("", "", "", "/..")
		}

		for _, entry := range entries {
			size := entry.GetSize()
			filling := int(size * 10 / largest)

			filled := getPounds(filling)
			notFilled := getPounds(filling - 10)

			ui.Table.AddRow(" ",
				humanize.Bytes(size),
				"["+filled+notFilled+"]",
				entry.GetName(),
			)
		}
		ui.Table.Select(0)
	}
}

func main() {
	var concurrent int
	var threads int
	var rootPath string

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flag.IntVar(&concurrent, "c", 256, "number of concurrent workers")
	flag.IntVar(&threads, "t", 4, "number of threads on which to spawn goroutines")
	flag.StringVar(&rootPath, "p", cwd, "path to start in")
	flag.Parse()

	runtime.GOMAXPROCS(threads)

	open := make(chan bool, concurrent)

	root = NewNode(".")
	RecursiveCompute(open, &root, rootPath)

	// ui
	ui := NewUI()
	done := make(chan int)
	defer ui.Close()

	ui.SelectHandler = func(selected string) {
		selected = strings.TrimPrefix(selected, "/")
		if selected == ".." {
			current = current.Parent
			load(&ui, current)
		} else {
			entry, ok := current.Get(selected)
			if ok && entry.Kind() == "Directory" {
				current = entry.(*TreeNode)
				load(&ui, current)
			}
		}
	}

	ui.EscHandler = func() {
		if current == &root {
			done <- 1
		} else {
			current = current.Parent
			load(&ui, current)
		}
	}

	go func() {
		err = ui.Run()
		if err != nil {
			log.Fatal(err)
		}
		done <- 1
	}()

	current = &root
	load(&ui, current)
	<-done
}
