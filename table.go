package gdu

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Table struct {
	Path        string
	buffer      [][]string
	selectedRow int
	view        int
}

func MakeTable() Table {
	return Table{
		buffer:      make([][]string, 0),
		selectedRow: 0,
		view:        0,
	}
}

func (t *Table) Clear() {
	t.buffer = make([][]string, 0)
}

func (t *Table) AddRow(row ...string) {
	t.buffer = append(t.buffer, row)
}

func (t *Table) Select(row int) {
	t.selectedRow = row
}

func (t *Table) MoveDown() {
	if t.selectedRow < len(t.buffer)-1 {
		t.selectedRow += 1
	}
}

func printLine(row int, s string, fg, bg termbox.Attribute, fill rune, full bool) {
	width, _ := termbox.Size()
	for i := 0; i < width; i++ {
		var r rune
		if i < len(s) {
			r = rune(s[i])
		} else {
			if !full {
				break
			}
			r = fill
		}
		termbox.SetCell(i, row, r, fg, bg)
	}
}

func printRow(row int, data []string, columns []int, fg, bg termbox.Attribute, fill rune, full bool) {
	width, _ := termbox.Size()
	var r rune
	var i int = 0
	for j, s := range data {
		b := i
		for ; i < b+columns[j]; i++ {
			if i-b < len(s) {
				r = rune(s[i-b])
			} else {
				r = ' '
			}
			termbox.SetCell(i, row, r, fg, bg)
		}
		termbox.SetCell(i, row, ' ', fg, bg)
		i++
		// for i < b+columns[j] {
		// 	if i-b < len(s) {
		// 		r = rune(s[i])
		// 	} else {
		// 		if !full {
		// 			break
		// 		}
		// 		r = ' '
		// 	}
		// 	termbox.SetCell(i, row, r, fg, bg)
		// 	i++
		// }
	}
	if full {
		for ; i < width; i++ {
			termbox.SetCell(i, row, fill, fg, bg)
		}
	}
}

func (t *Table) Draw() {
	_, height := termbox.Size()

	// header
	topString := fmt.Sprintf("gdu %s ~ Use the arrow keys to navigate", VERSION)
	printLine(0, topString, termbox.ColorBlack, termbox.ColorWhite, ' ', true)

	topString = fmt.Sprintf("--- %s ", t.Path)
	printLine(1, topString, termbox.ColorWhite, termbox.ColorBlack, '-', true)

	columns := make([]int, 4)
	for _, row := range t.buffer {
		for i, s := range row {
			columns[i] = max(columns[i], len(s))
		}
	}

	for i := t.view; i < min(len(t.buffer), t.view+height-2); i++ {
		printRow(i+2, t.buffer[i], columns, termbox.ColorWhite, termbox.ColorBlack, ' ', false)
	}
}