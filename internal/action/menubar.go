package action

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/display"
	"github.com/zyedidia/micro/v2/internal/screen"
)

type MenuBar struct {
	menuBar       *display.MenuBar
	menu          *Menu
	selectedIndex int
}

type topMenuDefinition struct {
	title          string
	menuDefinition *[]display.MenuDefinition
}

var topMenus []topMenuDefinition
var topMenuNames []string

func init() {
	topMenus = []topMenuDefinition{
		{
			title: "File",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "New",
					Shortcut:   "Ctrl+t",
					ActionName: "AddTab",
				},
				{
					Title:      "New Terminal...",
					ActionName: "command-edit:term ",
				},

				{
					Title:    "",
					Shortcut: "",
				},
				{
					Title:      "Open...",
					Shortcut:   "Ctrl+o",
					ActionName: "OpenFile",
				},
				{
					Title:      "Save",
					Shortcut:   "Ctrl+s",
					ActionName: "Save",
				},
				{
					Title:      "Save As...",
					ActionName: "SaveAs",
				},
				{
					Title:      "Close",
					Shortcut:   "Ctrl+q",
					ActionName: "Quit",
				},
				{
					Title:    "",
					Shortcut: "",
				},
				{
					Title:      "Quit",
					ActionName: "QuitAll",
				},
			},
		},
		{
			title: "Edit",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "Undo",
					Shortcut:   "Ctrl+z",
					ActionName: "Undo",
				},
				{
					Title:      "Redo",
					Shortcut:   "Ctrl+y",
					ActionName: "Redo",
				},
				{
					Title:    "",
					Shortcut: "",
				},
				{
					Title:      "Cut",
					Shortcut:   "Ctrl+x",
					ActionName: "Cut",
				},
				{
					Title:      "Copy",
					Shortcut:   "Ctrl+c",
					ActionName: "Copy",
				},
				{
					Title:      "Paste",
					Shortcut:   "Ctrl+v",
					ActionName: "Paste",
				},
				{
					Title:    "",
					Shortcut: "",
				},
				{
					Title:      "Find...",
					Shortcut:   "Ctrl+f",
					ActionName: "Find",
				},
				{
					Title:      "Find Next",
					Shortcut:   "Ctrl+n",
					ActionName: "FindNext",
				},
				{
					Title:      "Find Previous",
					Shortcut:   "Ctrl+p",
					ActionName: "FindPrevious",
				},
			},
		},
		{
			title: "Selection",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "Select All",
					ActionName: "SelectAll",
				},
				{
					Title:    "",
					Shortcut: "",
				},
				{
					Title:      "Spawn Cursors on Selection",
					Shortcut:   "Alt+m",
					ActionName: "SpawnMultiCursorSelect",
				},
				{
					Title:      "Spawn Cursor Up",
					Shortcut:   "Alt+Shift+Up",
					ActionName: "SpawnMultiCursorUp",
				},
				{
					Title:      "Spawn Cursor Down",
					Shortcut:   "Alt+Shift+Down",
					ActionName: "SpawnMultiCursorDown",
				},
				{
					Title:      "Remove Cursor",
					Shortcut:   "Alt-p",
					ActionName: "RemoveMultiCursor",
				},
				{
					Title:      "Remove All Cursors",
					Shortcut:   "Alt-c",
					ActionName: "RemoveAllMultiCursors",
				},
			},
		},
		{
			title: "Go",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "Go to Line...",
					Shortcut:   "Ctrl+l",
					ActionName: "command-edit:goto ",
				},
				{
					Title:      "Go to Next Paragraph",
					Shortcut:   "Alt+}",
					ActionName: "ParagraphNext",
				},
				{
					Title:      "Go to Previous Paragraph",
					Shortcut:   "Alt+{",
					ActionName: "ParagraphPrevious",
				},
				{
					Title:      "Go to Matching Bracket/Brace",
					ActionName: "JumpToMatchingBrace",
				},
				{
					Title:      "Go to Next Tab",
					Shortcut:   "Alt-.",
					ActionName: "NextTab|FirstTab",
				},
				{
					Title:      "Go to Previous Tab",
					Shortcut:   "Alt-,",
					ActionName: "PreviousTab|LastTab",
				},
				{
					Title:      "Go to Next Diff",
					Shortcut:   "Alt-]",
					ActionName: "DiffNext|CursorEnd",
				},
				{
					Title:      "Go to Previous Diff",
					Shortcut:   "Alt-[",
					ActionName: "DiffPrevious|CursorStart",
				},
			},
		},
		{
			title: "Help",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "Help",
					ActionName: "ToggleHelp",
				},
			},
		},
	}
	topMenuNames = []string{}
	for _, menu := range topMenus {
		topMenuNames = append(topMenuNames, menu.title)
	}
}

func NewMenuBar() *MenuBar {
	m := new(MenuBar)
	m.menuBar = display.NewMenuBar(topMenuNames)
	m.menu = NewMenu()
	return m
}

func (m *MenuBar) SetActionChan(actionChan chan string) {
	m.menu.SetActionChan(actionChan)
}

func (m *MenuBar) IsOpen() bool {
	return m.menu.IsOpen
}

func (m *MenuBar) Display() {
	m.menuBar.Display()
	if m.menu.IsOpen {
		m.menu.Display()
		screen.Screen.HideCursor()
	}
}

func (m *MenuBar) HandleEvent(event tcell.Event) {
	switch e := event.(type) {

	case *tcell.EventKey:
		if m.menu.IsOpen {
			switch e.Key() {
			case tcell.KeyLeft:
				m.selectedIndex--
				if m.selectedIndex < 0 {
					m.selectedIndex = 0
				}
				m.menu.Open(m.menuBar.MenuIndexLeft(m.selectedIndex), 1, topMenus[m.selectedIndex].menuDefinition)
				m.menuBar.SetSelected(m.selectedIndex)
			case tcell.KeyRight:
				m.selectedIndex++
				if m.selectedIndex >= len(topMenus) {
					m.selectedIndex = len(topMenus) - 1
				}
				m.menu.Open(m.menuBar.MenuIndexLeft(m.selectedIndex), 1, topMenus[m.selectedIndex].menuDefinition)
				m.menuBar.SetSelected(m.selectedIndex)
			}
		}

	case *tcell.EventMouse:
		mx, my := e.Position()
		if my != 0 {
			break
		}
		switch e.Buttons() {
		case tcell.Button1:
			selectedMenuIndex, leftX := m.menuBar.MenuItemIndexAtX(mx)
			if selectedMenuIndex == m.selectedIndex && m.menu.IsOpen {
				m.menu.IsOpen = false
				m.menuBar.SetSelected(-1)
			} else if selectedMenuIndex != -1 {
				m.menu.Open(leftX, 1, topMenus[selectedMenuIndex].menuDefinition)
				m.menuBar.SetSelected(selectedMenuIndex)
				m.selectedIndex = selectedMenuIndex
			}
		}
		return
	}

	m.menu.HandleEvent(event)
	if !m.menu.IsOpen {
		m.menuBar.SetSelected(-1)
	}
}
