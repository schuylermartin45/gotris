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
}

// RenderGame runs the primary gameplay loop.
func (t *TextGame) RenderGame() {
	// A digital frontier...
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
}

/*
 Draws the current board to the screen.
*/
func (t TextGame) drawBoard() {
	workingGrid := t.board.Current()
	for row := 0; row < len(workingGrid); row++ {
		var mask uint8 = 1 << 7
		for col := 0; col < int(model.BoardWidth); col++ {
			if (workingGrid[row] & mask) > 0 {
				// TODO draw piece
				t.screen.SetContent(2*col, row, '█', nil, 32)
				t.screen.SetContent((2*col)+1, row, '█', nil, 32)
			} else {
				// TODO draw transparency/background/nothing
				t.screen.SetContent(2*col, row, 32, nil, 32)
				t.screen.SetContent((2*col)+1, row, 32, nil, 32)
			}
			mask >>= 1
		}
	}
	t.screen.Show()
}
