package action

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/display"
)

type MenuBar struct {
	menuBar *display.MenuBar
	menu    *Menu
}

var fileMenuDefinition *[]display.MenuDefinition

func init() {
	fileMenuDefinition = &[]display.MenuDefinition{
		{
			Title:    "Open",
			Shortcut: "Ctrl+o",
		},
		{
			Title:    "Save",
			Shortcut: "Ctrl+s",
		},
		{
			Title:    "Close",
			Shortcut: "Ctrl+w",
		},
	}
}

func NewMenuBar(menu *Menu) *MenuBar {
	m := new(MenuBar)
	m.menuBar = display.NewMenuBar()
	m.menu = menu
	return m
}

func (m *MenuBar) Display() {
	m.menuBar.Display()
}

func (m *MenuBar) HandleEvent(event tcell.Event) {
	switch e := event.(type) {
	case *tcell.EventMouse:
		_, my := e.Position()
		if my != 0 {
			return
		}
		switch e.Buttons() {
		case tcell.Button1:
			m.menu.Open(0, 1, fileMenuDefinition)
		}
	}
}
