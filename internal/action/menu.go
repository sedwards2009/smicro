package action

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/display"
)

type Menu struct {
	menuWindow *display.MenuWindow
	IsOpen     bool
}

func NewMenu() *Menu {
	d := new(Menu)
	d.menuWindow = display.NewMenuWindow()
	d.IsOpen = false
	return d
}

func (d *Menu) Open(x int, y int, menuDefinition *[]display.MenuDefinition) {
	d.IsOpen = true
	d.menuWindow.X = x
	d.menuWindow.Y = y
	d.menuWindow.MenuDefinition = menuDefinition
	d.menuWindow.SelectedRow = 0
}

func (d *Menu) Display() {
	if !d.IsOpen {
		return
	}

	d.menuWindow.Display()
}

func (d *Menu) HandleEvent(event tcell.Event) {
	if !d.IsOpen {
		return
	}

	switch e := event.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEsc:
			d.IsOpen = false
		case tcell.KeyUp:
			d.menuWindow.SelectedRow = d.nextMenuItem(-1)
		case tcell.KeyDown:
			d.menuWindow.SelectedRow = d.nextMenuItem(1)
		}
	}
}

func (d *Menu) nextMenuItem(direction int) int {
	selectedRow := d.menuWindow.SelectedRow
	menuDefinition := d.menuWindow.MenuDefinition

	next := func() {
		selectedRow += direction

		if selectedRow < 0 {
			selectedRow = len(*menuDefinition) - 1
		}
		if selectedRow >= len(*menuDefinition) {
			selectedRow = 0
		}
	}
	next()

	for (*menuDefinition)[selectedRow].Title == "" {
		next()
	}
	return selectedRow
}
