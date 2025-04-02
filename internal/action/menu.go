package action

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/display"
)

type Menu struct {
	menuWindow *display.MenuWindow
	IsOpen     bool
	actionChan chan string
}

func NewMenu() *Menu {
	d := new(Menu)
	d.menuWindow = display.NewMenuWindow()
	d.IsOpen = false
	return d
}

func (m *Menu) SetActionChan(actionChan chan string) {
	m.actionChan = actionChan
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
		case tcell.KeyEnter:
			row := d.menuWindow.SelectedRow
			d.executeMenuItem(row)
		}

	case *tcell.EventMouse:
		if e.Buttons() == tcell.Button1 {
			x, y, w, h := d.menuWindow.Position()
			mx, my := e.Position()
			// Check if the mouse is within the menu window
			inside := mx >= x && mx < x+w && my >= y && my < y+h
			if !inside {
				d.IsOpen = false
			} else {
				// Execute the menu item
				row := my - y - 1
				if row >= 0 && row < len(*d.menuWindow.MenuDefinition) {
					d.executeMenuItem(row)
				}
			}
		}
	}
}

func (d *Menu) executeMenuItem(row int) {
	menuDefinition := (*d.menuWindow.MenuDefinition)[row]
	if menuDefinition.Title != "" {
		actionName := menuDefinition.ActionName
		if actionName != "" {
			d.actionChan <- actionName
		}
		d.IsOpen = false
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
