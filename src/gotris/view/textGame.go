/*
 * File:        textGame.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: An advanced gameplay mode that runs in a text terminal.
 */
package view

import (
	"../model"
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	"time"
)

/***** Types *****/

// TextGame renders Gotris in an interactive text-based UI.
type TextGame struct {
	board  *model.Board
	screen tcell.Screen
}

// Text Mode Color Enum
type color uint8

// Text Mode Color enumerations
const (
	Blue            color = 0
	Cyan            color = 1
	Grey            color = 2
	Yellow          color = 3
	Green           color = 4
	Violet          color = 5
	Red             color = 6
	BoardBackground color = 7
	BoardBorder     color = 8
	BoardForeground color = 9
	TextColor       color = 10
)

/***** Functions *****/

// lookupColor returns the `tcell` color code for a given color
func lookupColor(clr color) tcell.Style {
	bkgrd := tcell.ColorBlack
	frgrd := tcell.ColorDarkGrey
	style := tcell.StyleDefault
	switch clr {
	case Blue:
		return style.Foreground(tcell.ColorBlue).Background(tcell.ColorDarkBlue)
	case Cyan:
		return style.Foreground(tcell.ColorLightBlue).Background(tcell.ColorRoyalBlue)
	case Grey:
		return style.Foreground(tcell.ColorGrey).Background(tcell.ColorDarkGrey)
	case Yellow:
		return style.Foreground(tcell.ColorYellow).Background(tcell.ColorSandyBrown)
	case Green:
		return style.Foreground(tcell.ColorGreen).Background(tcell.ColorDarkGreen)
	case Violet:
		return style.Foreground(tcell.ColorViolet).Background(tcell.ColorDarkViolet)
	case Red:
		return style.Foreground(tcell.ColorRed).Background(tcell.ColorDarkRed)
	case TextColor:
		fallthrough
	case BoardBorder:
		return style.Foreground(frgrd).Background(bkgrd)
	case BoardForeground:
		return style.Foreground(frgrd).Background(tcell.ColorLightSlateGray)
	case BoardBackground:
		fallthrough
	default:
		return style.Background(bkgrd)
	}
}

// lookupTileColor maps TileColor to the `tcell` color code
func lookupTileColor(clr model.TileColor) tcell.Style {
	switch clr {
	case model.Blue:
		return lookupColor(Blue)
	case model.Cyan:
		return lookupColor(Cyan)
	case model.Grey:
		return lookupColor(Grey)
	case model.Yellow:
		return lookupColor(Yellow)
	case model.Green:
		return lookupColor(Green)
	case model.Violet:
		return lookupColor(Violet)
	case model.Red:
		return lookupColor(Red)
	}
	return 0x00
}

/***** Methods *****/

// RenderHelpMenu returns a string to display the help menu in the terminal.
func (t TextGame) RenderHelpMenu() string {
	return "Text Mode\n" +
		"\nAbout\n" +
		"  This mode is an advanced, real-time text-based gameplay mode.\n" +
		"  It is written using the `tcell` Go package.\n" +
		"\nControls\n" +
		"  * W/[Up]:         Rotate\n" +
		"  * A/[Left]:       Move left\n" +
		"  * S/[Down]:       Move right\n" +
		"  * D/[Right]:      Move down\n" +
		"  * [Space]:        Drop tile to floor\n" +
		"  * [Esc]/[Ctrl-C]: Exit game\n"
}

// InitGame initializes the game.
func (t *TextGame) InitGame(b *model.Board) {
	t.board = b

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var error error
	t.screen, error = tcell.NewScreen()
	if error != nil {
		fmt.Fprintf(os.Stderr, "%v\n", error)
		os.Exit(ERROR_SCREEN_INIT)
	}
	if error = t.screen.Init(); error != nil {
		fmt.Fprintf(os.Stderr, "%v\n", error)
		os.Exit(ERROR_SCREEN_INIT)
	}
	// Kick off event listener thread.
	go t.initEventListener()
}

// RenderGame runs the primary gameplay loop.
func (t *TextGame) RenderGame() {
	// Primary game loop loops until the game completes
	for {
		// Advance the game
		_, endGame := t.board.Next()
		t.drawBoard()

		// Draw the game. Game speed increases with level until a certain point.
		delay := 500 - int(50*t.board.GetLevel())
		if delay < 100 {
			delay = 100
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)

		// Stop the loop on the event that the game has ended.
		if endGame {
			break
		}
	}
}

// ExitGame is a callback triggered when the game terminates
func (t *TextGame) ExitGame(playAgain bool) {
	// Clean up screen object
	t.screen.Fini()
}

/*
 Draws a string.

 @param x   Left-top corner x position of the string
 @param y   Left-top corner y position of the string
 @param str String to draw
*/
func (t *TextGame) drawStr(x int, y int, str string) {
	sizeX, sizeY := t.screen.Size()
	if (x < 0) || (y < 0) || (y > sizeY) {
		return
	}
	for row := 0; row < len(str); row++ {
		screenX := x + row
		if screenX > sizeX {
			break
		}
		t.screen.SetContent(screenX, y, rune(str[row]), nil, lookupColor(TextColor))
	}
}

/*
 Draws the current board to the screen.
*/
func (t *TextGame) drawBoard() {
	const (
		// Starting coordinates for the board
		boardX = 10
		boardY = 4
		// Starting coordinates for the next tile preview (relative to the board)
		previewX = boardY + (2 * 8) + 8
		previewY = boardY + 2
		// Starting coordinates for the score (relative to the board)
		scoreX = previewX + 4
		scoreY = boardY
	)
	t.screen.Fill(' ', lookupColor(BoardBackground))

	// Draw the main board
	workingGrid := t.board.Current()
	y := boardY
	for row := 0; row < len(workingGrid); row++ {
		// `y` includes padding from the top of the screen
		var mask uint8 = 1 << 7
		for col := 0; col < int(model.BoardWidth); col++ {
			// Calculate the left and right block x coordinates
			xL := boardX + (2 * col)
			xR := boardX + (2 * col) + 1
			if (workingGrid[row] & mask) > 0 {
				t.screen.SetContent(xL, y, '▇', nil, lookupColor(Red))
				t.screen.SetContent(xR, y, '▇', nil, lookupColor(Red))
			} else {
				t.screen.SetContent(xL, y, ' ', nil, lookupColor(BoardForeground))
				t.screen.SetContent(xR, y, '.', nil, lookupColor(BoardForeground))
			}
			mask >>= 1
		}
		y++
	}

	// Draw the score
	t.drawStr(scoreX, scoreY, t.board.GetDisplayScore())

	// Draw the next tile
	nextTile := t.board.GetNextTile()
	nextTileBlock := nextTile.GetBlock()
	nextTileColor := nextTile.GetColor()
	y = previewY
	// Pad the top of the preview view
	for col := 4; col < (2*model.TileSize)+4; col++ {
		x := previewX + col
		t.screen.SetContent(x, y, ' ', nil, lookupColor(BoardForeground))
	}
	y++
	// We draw outside of the rendering of the tile to provide lower padding.
	for row := 0; row < model.TileSize+1; row++ {
		// To save on rendering time, skip the first 2 columns, which we know are
		// padded to be empty on the initial tile's shape/orientation.
		var mask uint8 = 1 << 5
		for col := 2; col < model.TileSize+2; col++ {
			xL := previewX + (2 * col)
			xR := previewX + (2 * col) + 1
			// The check against row provides lower padding to the preview.
			if (row < model.TileSize) && ((nextTileBlock[row] & mask) > 0) {
				t.screen.SetContent(xL, y, '▇', nil, lookupTileColor(nextTileColor))
				t.screen.SetContent(xR, y, '▇', nil, lookupTileColor(nextTileColor))
			} else {
				t.screen.SetContent(xL, y, ' ', nil, lookupColor(BoardForeground))
				t.screen.SetContent(xR, y, ' ', nil, lookupColor(BoardForeground))
			}
			mask >>= 1
		}
		y++
	}

	// Render it all
	t.screen.Show()
}

/*
 Initializes the event listener
*/
func (t *TextGame) initEventListener() {
	for {
		event := t.screen.PollEvent()
		switch eventType := event.(type) {
		case *tcell.EventKey:
			var action Action = ActionIllegal
			switch eventType.Key() {
			// ASCII keys have to be handled separately
			case tcell.KeyRune:
				switch eventType.Rune() {
				case 'a':
					action = ActionLeft
				case 'd':
					action = ActionRight
				// Down
				case 's':
					action = ActionDown
				// Fast Down
				case 'w':
					action = ActionRotate
				case ' ':
					action = ActionFastDown
				}
			case tcell.KeyLeft:
				action = ActionLeft
			case tcell.KeyRight:
				action = ActionRight
			case tcell.KeyDown:
				action = ActionDown
			case tcell.KeyUp:
				action = ActionRotate
			// Exit
			case tcell.KeyCtrlC:
				fallthrough
			case tcell.KeyEsc:
				action = ActionExit
			}
			if action != ActionIllegal {
				ActionHandler(t.board, action, func() {
					t.screen.Fini()
					os.Exit(EXIT_SUCCESS)
					return
				})
				// Re-render the board on action to make visual feedback more
				// apparent
				t.drawBoard()
			}
		default:
			continue
		}
	}
}
