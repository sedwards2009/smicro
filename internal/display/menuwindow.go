package display

import (
	runewidth "github.com/mattn/go-runewidth"
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/screen"
)

type MenuDefinition struct {
	Title      string
	Shortcut   string
	ActionName string
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

func (d *MenuWindow) MoveOnScreen() {
	if d.X < 0 {
		d.X = 0
	}
	if d.Y < 0 {
		d.Y = 0
	}

	menuWidth := d.menuWidthInCells(d.MenuDefinition)
	width, height := screen.Screen.Size()
	if d.X+menuWidth > width {
		d.X = width - menuWidth
	}
	if d.Y+len(*d.MenuDefinition)+2 > height {
		d.Y = height - len(*d.MenuDefinition) + 2
	}
}

func (d *MenuWindow) menuWidthInCells(md *[]MenuDefinition) int {
	titleWidth, shortcutWidth := measureWidths(md)
	return 1 + 1 + titleWidth + 2 + shortcutWidth + 1 + 1
}

func (d *MenuWindow) Display() {
	titleWidth, _ := measureWidths(d.MenuDefinition)
	y := d.Y
	menuWidth := d.menuWidthInCells(d.MenuDefinition)

	borderStyle := config.Colorscheme["tabbar"]

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

	// Draw the drop shadow
	drawDimVerticalLine(d.X+menuWidth, d.Y+1, len(*d.MenuDefinition)+1)
	drawDimVerticalLine(d.X+menuWidth+1, d.Y+1, len(*d.MenuDefinition)+1)
	drawDimHorizontalLine(d.X+2, y+1, menuWidth)
}

func measureWidths(menuDefinition *[]MenuDefinition) (int, int) {
	maxTitleWidth := 0
	maxShortcutWidth := 0
	for _, item := range *menuDefinition {
		width := runewidth.StringWidth(item.Title)
		if width > maxTitleWidth {
			maxTitleWidth = width
		}
		shortCutWidth := runewidth.StringWidth(item.Shortcut)
		if shortCutWidth > maxShortcutWidth {
			maxShortcutWidth = shortCutWidth
		}
	}
	return maxTitleWidth, maxShortcutWidth
}

func (d *MenuWindow) Position() (x int, y int, w int, h int) {
	titleWidth, shortcutWidth := measureWidths(d.MenuDefinition)
	totalWidth := 2 + titleWidth + 2 + shortcutWidth + 2
	return d.X, d.Y, totalWidth, 2 + len(*d.MenuDefinition)
}

func dimCell(x int, y int) {
	cellRune, cellRunes, style := screen.GetContent(x, y)
	fg, bg, _ := style.Decompose()

	fgR, fgG, fgB := fg.TrueColor().RGB()
	fg = tcell.NewRGBColor(fgR/2, fgG/2, fgB/2)

	bgR, bgG, bgB := bg.TrueColor().RGB()
	bg = tcell.NewRGBColor(bgR/2, bgG/2, bgB/2)

	style = style.Foreground(fg).Background(bg)
	screen.SetContent(x, y, cellRune, cellRunes, style)
}

func drawDimHorizontalLine(x int, y int, length int) {
	for i := 0; i < length; i++ {
		dimCell(x+i, y)
	}
}

func drawDimVerticalLine(x int, y int, length int) {
	for i := 0; i < length; i++ {
		dimCell(x, y+i)
	}
}
