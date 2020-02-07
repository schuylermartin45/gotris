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
)

/***** Functions *****/

// lookupColor returns the `tcell` color code for a given color
func lookupColor(clr color) tcell.Style {
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
	case BoardBackground:
		return style.Background(tcell.ColorBlack)
	case BoardBorder:
		return style.Background(tcell.ColorDarkGrey)
	case BoardForeground:
		return style.Background(tcell.ColorLightSlateGray)
	}
	return 0x00
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
		// TODO rm this line
		"  NOTE: this game mode is not complete or playable yet.\n" +
		"\nControls\n" +
		"  * W:       Drop tile to floor\n" +
		"  * A:       Move left\n" +
		"  * S:       Move right\n" +
		"  * D:       Move down\n" +
		"  * [Space]: Rotate\n" +
		"  * E:       Exit game\n"
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

		// Draw the game
		time.Sleep(100 * time.Millisecond)

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
 Draws the current board to the screen.
*/
func (t TextGame) drawBoard() {
	const (
		// Starting coordinates for the board
		boardX = 8
		boardY = 2
		// Starting coordinates for the next tile preview (relative to the board)
		previewX = boardY + (2 * 8) + 8
		previewY = boardY
	)

	t.screen.Fill(' ', lookupColor(BoardBackground))

	// Draw the board outline

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
				t.screen.SetContent(xR, y, ' ', nil, lookupColor(BoardForeground))
			}
			mask >>= 1
		}
		y++
	}

	// Draw the score

	// Draw the next tile
	nextTile := t.board.GetNextTile()
	nextTileBlock := nextTile.GetBlock()
	nextTileColor := nextTile.GetColor()
	y = previewY
	for row := 0; row < model.TileSize; row++ {
		// To save on rendering time, skip the first 2 columns, which we know are
		// padded to be empty on the initial tile's shape/orientation.
		var mask uint8 = 1 << 5
		for col := 2; col < model.TileSize+2; col++ {
			xL := previewX + (2 * col)
			xR := previewX + (2 * col) + 1
			if (nextTileBlock[row] & mask) > 0 {
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
 Action handler. Given an action, performs a board operation.

 @param action	Action to interpret
*/
func (t *TextGame) actionHandler(action Action) {
	switch action {
	case ActionLeft:
		t.board.MoveLeft()
	case ActionRight:
		t.board.MoveRight()
	case ActionDown:
		t.board.MoveDown()
	case ActionFastDown:
		t.board.MoveFastDown()
	case ActionRotate:
		t.board.Rotate()
	case ActionExit:
		t.screen.Fini()
		os.Exit(EXIT_SUCCESS)
	}
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
					action = ActionFastDown
				case ' ':
					action = ActionRotate
				}
			case tcell.KeyLeft:
				action = ActionLeft
			case tcell.KeyRight:
				action = ActionRight
			case tcell.KeyDown:
				action = ActionDown
			case tcell.KeyUp:
				action = ActionFastDown
			// Exit
			case tcell.KeyCtrlC:
				fallthrough
			case tcell.KeyEsc:
				action = ActionExit
			}
			if action != ActionIllegal {
				t.actionHandler(action)
			}
		default:
			continue
		}
	}
}
