package display

import (
	runewidth "github.com/mattn/go-runewidth"
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/screen"
)

type MenuDefinition struct {
	Title    string
	Shortcut string
}

type MenuWindow struct {
	MenuDefinition *[]MenuDefinition
	X              int
	Y              int
	SelectedRow    int
}

func NewMenuWindow() *MenuWindow {
	return new(MenuWindow)
}

func drawHorizontalLine(x int, y int, width int, borderStyle tcell.Style, middleStyle tcell.Style, left rune,
	middle rune, right rune) {

	screen.SetContent(x, y, left, nil, borderStyle)
	for i := 1; i < width-1; i++ {
		screen.SetContent(x+i, y, middle, nil, middleStyle)
	}
	screen.SetContent(x+width-1, y, right, nil, borderStyle)
}

func drawText(x int, y int, text string, style tcell.Style) {
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

func (d *MenuWindow) Display() {
	titleWidth, shortcutWidth := measureWidths(d.MenuDefinition)
	y := d.Y
	menuWidth := 1 + 1 + titleWidth + 2 + shortcutWidth + 1 + 1

	backgroundColor := tcell.ColorNavy
	borderStyle := tcell.Style{}.Foreground(tcell.ColorWhite).Background(backgroundColor).Bold(true)

	drawHorizontalLine(d.X, y, menuWidth, borderStyle, borderStyle, '┌', '─', '┐')
	y++

	for i, item := range *d.MenuDefinition {
		if item.Title == "" {
			drawHorizontalLine(d.X, y, menuWidth, borderStyle, borderStyle, '├', '─', '┤')
		} else {
			textStyle := borderStyle
			if i == d.SelectedRow {
				textStyle = textStyle.Reverse(true)
			}
			drawHorizontalLine(d.X, y, menuWidth, borderStyle, textStyle, '│', ' ', '│')
			drawText(d.X+2, y, item.Title, textStyle)
			drawText(d.X+2+titleWidth+2, y, item.Shortcut, textStyle)
		}
		y++
	}

	drawHorizontalLine(d.X, y, menuWidth, borderStyle, borderStyle, '└', '─', '┘')
}

func measureWidths(menuDefinition *[]MenuDefinition) (int, int) {
	maxWidth := 0
	maxShortcutWidth := 0
	for _, item := range *menuDefinition {
		width := runewidth.StringWidth(item.Title)
		if width > maxWidth {
			maxWidth = width
		}
		shortCutWidth := runewidth.StringWidth(item.Shortcut)
		if shortCutWidth > maxShortcutWidth {
			maxShortcutWidth = shortCutWidth
		}
	}
	return maxWidth, maxShortcutWidth
}
