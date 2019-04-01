package main

import (
	"flag"
	"log"
	"os"
	"sort"
	"strings"

	"gdu"

	"github.com/dustin/go-humanize"
)

var err error
var current *gdu.TreeNode

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

func load(ui *gdu.UI, root *gdu.TreeNode) {
	ui.Table.Clear()
	ui.Table.Path = root.Path
	entries := root.GetEntries()

	if len(entries) > 0 {
		sort.Sort(entries)

		// now entries[0] is guaranteed to be the largest
		largest := entries[0].GetSize()

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
	flag.IntVar(&concurrent, "c", 940, "number of concurrent workers")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	open := make(chan bool, concurrent)

	root := gdu.NewNode(".")
	gdu.RecursiveCompute(open, &root, cwd)

	// ui

	ui := gdu.NewUI()
	done := make(chan int)
	defer ui.Close()

	ui.SelectHandler = func(selected string) {
		entry, ok := current.Get(strings.TrimPrefix(selected, "/"))
		if ok && entry.Kind() == "Directory" {
			current = entry.(*gdu.TreeNode)
			load(&ui, current)
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
