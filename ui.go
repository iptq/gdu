package gdu

import (
	"github.com/nsf/termbox-go"
)

type UI struct {
	Table         Table
	Events        chan termbox.Event
	Redraw        chan int
	EscHandler    func()
	SelectHandler func(string)
}

func NewUI() UI {
	return UI{
		Table:  MakeTable(),
		Events: make(chan termbox.Event),
		Redraw: make(chan int),
	}
}

func (ui *UI) eventLoop() {
	for {
		event := termbox.PollEvent()
		ui.Events <- event
	}
}

func (ui *UI) Close() {
	termbox.Close()
}

func (ui *UI) Run() (err error) {
	err = termbox.Init()
	if err != nil {
		return
	}

	go ui.eventLoop()

	termbox.SetCursor(-1, -1)
	for {
		ui.Table.Draw()
		termbox.Flush()

		select {
		case event := <-ui.Events:
			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyArrowUp:
					ui.Table.MoveUp()

				case termbox.KeyArrowDown:
					ui.Table.MoveDown()

				case termbox.KeyEnter:
					ui.SelectHandler(ui.Table.GetSelectedName())

				case termbox.KeyEsc:
					ui.EscHandler()

				case termbox.KeyCtrlC:
					return
				}
			}
		case _ = <-ui.Redraw:
			break
		}
	}
}
