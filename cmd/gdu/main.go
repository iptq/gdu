package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"gdu"

	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var err error
var current *gdu.TreeNode

func load(table *tview.Table, root *gdu.TreeNode) {
	table.Clear()
	entries := root.GetEntries()
	sort.Sort(sort.Reverse(entries))
	for i, entry := range entries {
		table.SetCell(i, 0,
			tview.NewTableCell(entry.GetName()))
		table.SetCell(i, 1,
			tview.NewTableCell(humanize.Bytes(entry.GetSize())))
	}
	table.Select(0, 0)
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	root := gdu.NewNode(".")
	gdu.RecursiveCompute(&root, cwd)
	fmt.Println("waiting...")

	// ui

	app := tview.NewApplication()

	current = &root
	table := tview.NewTable().
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical)
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

		entry := current.Get(name)
		if entry.Kind() == "Directory" {
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
