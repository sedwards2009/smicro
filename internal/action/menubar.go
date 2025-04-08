package action

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/config"
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

func NewMenuBar() *MenuBar {
	m := new(MenuBar)
	m.menuBar = display.NewMenuBar(topMenuNames)
	m.menu = NewMenu()
	return m
}

func init() {
	topMenus = []topMenuDefinition{
		{
			title: "File",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "New",
					ActionName: "AddTab",
				},
				{
					Title:      "New Terminal...",
					ActionName: "command-edit:term ",
				},
				{
					Title: "",
				},
				{
					Title:      "Open...",
					ActionName: "OpenFile",
				},
				{
					Title:      "Save",
					ActionName: "Save",
				},
				{
					Title:      "Save As...",
					ActionName: "SaveAs",
				},
				{
					Title:      "Close",
					ActionName: "Quit",
				},
			},
		},
		{
			title: "Edit",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "Undo",
					ActionName: "Undo",
				},
				{
					Title:      "Redo",
					ActionName: "Redo",
				},
				{
					Title: "",
				},
				{
					Title:      "Cut",
					ActionName: "Cut|CutLine",
				},
				{
					Title:      "Copy",
					ActionName: "Copy|CopyLine",
				},
				{
					Title:      "Paste",
					ActionName: "Paste",
				},
				{
					Title: "",
				},
				{
					Title:      "Find...",
					ActionName: "Find",
				},
				{
					Title:      "Find Next",
					ActionName: "FindNext",
				},
				{
					Title:      "Find Previous",
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
					ActionName: "SpawnMultiCursorSelect",
				},
				{
					Title:      "Spawn Cursor Up",
					ActionName: "SpawnMultiCursorUp",
				},
				{
					Title:      "Spawn Cursor Down",
					ActionName: "SpawnMultiCursorDown",
				},
				{
					Title:      "Remove Cursor",
					ActionName: "RemoveMultiCursor",
				},
				{
					Title:      "Remove All Cursors",
					ActionName: "RemoveAllMultiCursors",
				},
			},
		},
		{
			title: "Go",
			menuDefinition: &[]display.MenuDefinition{
				{
					Title:      "Go to Line...",
					ActionName: "command-edit:goto ",
				},
				{
					Title:      "Jump...",
					ActionName: "command-edit:jump ",
				},
				{
					Title:      "Go to Next Paragraph",
					ActionName: "ParagraphNext",
				},
				{
					Title:      "Go to Previous Paragraph",
					ActionName: "ParagraphPrevious",
				},
				{
					Title:      "Go to Matching Bracket/Brace",
					ActionName: "JumpToMatchingBrace",
				},
				{
					Title:      "Go to Next Tab",
					ActionName: "NextTab|FirstTab",
				},
				{
					Title:      "Go to Previous Tab",
					ActionName: "PreviousTab|LastTab",
				},
				{
					Title:      "Go to Next Diff",
					ActionName: "DiffNext|CursorEnd",
				},
				{
					Title:      "Go to Previous Diff",
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
				{
					Title:      "Colors",
					ActionName: "command:help colors",
				},
				{
					Title:      "Default keys",
					ActionName: "command:help defaultkeys",
				},
				{
					Title:      "Options",
					ActionName: "command:help options",
				},
				{
					Title:      "Commands",
					ActionName: "command:help commands",
				},
				{
					Title:      "Plugins",
					ActionName: "command:help plugins",
				},
				{
					Title:      "Copy & Paste",
					ActionName: "command:help copypaste",
				},
				{
					Title:      "Keybindings",
					ActionName: "command:help keybindings",
				},
				{
					Title:      "Tutorial",
					ActionName: "command:help tutorial",
				},
			},
		},
	}
	topMenuNames = []string{}
	for _, menu := range topMenus {
		topMenuNames = append(topMenuNames, menu.title)
	}
}

func (m *MenuBar) InitBindings() {
	for i := range topMenus {
		menu := &topMenus[i]
		for j := range *menu.menuDefinition {
			item := &(*menu.menuDefinition)[j]
			if item.ActionName != "" {
				item.Shortcut = actionToKeyBinding(item.ActionName)
			}
		}
	}
}

func actionToKeyBinding(actionName string) string {
	result := ""
	for _, category := range config.Bindings {
		for key, action := range category {
			if action == actionName {
				result = key
				if result[0] != 'F' { // Bias the search to favour non function key bindings
					return key
				}
			}
		}
	}
	return result
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
