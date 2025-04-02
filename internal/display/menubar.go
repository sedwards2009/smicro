package display

import (
	runewidth "github.com/mattn/go-runewidth"
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/screen"
)

type MenuBar struct {
	menuNames []string
	selected  int
}

var menuItemPadding = 1
var menuItemSpacing = 3

func NewMenuBar(menuNames []string) *MenuBar {
	return &MenuBar{
		menuNames: menuNames,
		selected:  -1,
	}
}

func (m *MenuBar) Display() {
	tabBarColor := tcell.ColorNavy
	tabBarStyle := tcell.Style{}.Foreground(tcell.ColorWhite).Background(tabBarColor)

	w, _ := screen.Screen.Size()

	for i := 0; i < w; i++ {
		screen.SetContent(i, 0, ' ', nil, tabBarStyle)
	}

	x := 1
	drawPadding := func(style tcell.Style) {
		for i := 0; i < menuItemPadding; i++ {
			screen.SetContent(x+i, 0, ' ', nil, style)
		}
	}

	for i, name := range m.menuNames {
		menuItemStyle := tabBarStyle
		if i == m.selected {
			menuItemStyle = menuItemStyle.Reverse(true)
		}

		drawPadding(menuItemStyle)
		x += menuItemPadding

		drawText(x, 0, name, menuItemStyle)
		x += runewidth.StringWidth(name)

		drawPadding(menuItemStyle)
		x += menuItemPadding

		x += menuItemSpacing
	}
}

func (m *MenuBar) MenuItemIndexAtX(posX int) (index int, leftX int) {
	x := 1
	for i, name := range m.menuNames {
		if posX < x {
			return -1, -1
		}

		left := x
		x += menuItemPadding
		x += runewidth.StringWidth(name)
		x += menuItemPadding
		if posX < x {
			return i, left
		}

		x += menuItemSpacing
	}
	return -1, -1
}

func (m *MenuBar) MenuIndexLeft(index int) int {
	x := 1
	for i, name := range m.menuNames {
		if i == index {
			return x
		}

		x += menuItemPadding
		x += runewidth.StringWidth(name)
		x += menuItemPadding
		x += menuItemSpacing
	}
	return -1
}

func (m *MenuBar) SetSelected(index int) {
	m.selected = index
}
