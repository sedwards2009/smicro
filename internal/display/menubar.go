package display

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/screen"
)

type MenuBar struct {
}

func NewMenuBar() *MenuBar {
	return new(MenuBar)
}

func (m *MenuBar) Display() {
	tabBarColor := tcell.ColorNavy
	tabBarStyle := tcell.Style{}.Foreground(tcell.ColorBlack).Background(tabBarColor).Bold(true)

	w, _ := screen.Screen.Size()

	for i := 0; i < w; i++ {
		screen.SetContent(i, 0, ' ', nil, tabBarStyle)
	}
}
