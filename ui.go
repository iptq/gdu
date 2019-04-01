package gdu

import (
	"github.com/nsf/termbox-go"
)

type UI struct {
	Table         Table
	Events        chan termbox.Event
	Redraw        chan int
	SelectHandler func(int)
}

func NewUI() UI {
	return UI{
		Table:  MakeTable(),
		Events: make(chan termbox.Event),
		Redraw: make(chan int),
	}
}

func (ui UI) eventLoop() {
	for {
		event := termbox.PollEvent()
		ui.Events <- event
	}
}

func (ui UI) Run() (err error) {
	err = termbox.Init()
	if err != nil {
		return
	}
	defer termbox.Close()

	go ui.eventLoop()

	for {
		ui.Table.Draw()
		termbox.SetCursor(-1, -1)
		termbox.Flush()

		select {
		case event := <-ui.Events:
			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyArrowDown:
					ui.Table.MoveDown()
				case termbox.KeyCtrlC:
					return
				}
			}
		case _ = <-ui.Redraw:
		}
	}
}
