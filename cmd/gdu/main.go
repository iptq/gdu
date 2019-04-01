package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gdu"

	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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

func load(table *tview.Table, root *gdu.TreeNode) {
	table.Clear()
	entries := root.GetEntries()
	if len(entries) > 0 {
		sort.Sort(entries)

		// now entries[0] is guaranteed to be the largest
		largest := entries[0].GetSize()

		for i, entry := range entries {
			size := entry.GetSize()
			filling := int(size * 10 / largest)

			filled := getPounds(filling)
			notFilled := getPounds(filling - 10)

			table.SetCell(i, 0, tview.NewTableCell(" "))
			table.SetCell(i, 1, tview.NewTableCell(humanize.Bytes(size)).SetAlign(tview.AlignRight))
			table.SetCell(i, 2, tview.NewTableCell("["+filled+notFilled+"]"))
			table.SetCell(i, 3, tview.NewTableCell(entry.GetName()))
		}
		table.Select(0, 0)
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
	fmt.Println("waiting...")
	gdu.RecursiveCompute(open, &root, cwd)

	// ui

	app := tview.NewApplication()

	current = &root
	table := tview.NewTable().
		SetSelectable(true, false).
		SetSeparator(' ')
	load(table, &root)

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			if current == &root {
				app.Stop()
			} else {
				current = current.Parent
				load(table, current)
			}
		}
	}).SetSelectedFunc(func(r, c int) {
		cell := table.GetCell(r, c)
		name := cell.Text

		entry, ok := current.Get(name)
		if ok && entry.Kind() == "Directory" {
			current = entry.(*gdu.TreeNode)
			load(table, current)
		}
	})
	if err := app.SetRoot(table, true).
		SetFocus(table).
		Run(); err != nil {
		panic(err)
	}
}
