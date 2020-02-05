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
	board *model.Board
}

/***** Methods *****/

// InitGame initializes the game.
func (t TextGame) InitGame(b *model.Board) {
	t.board = b
}

// RenderGame runs the primary gameplay loop.
func (t TextGame) RenderGame() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, error := tcell.NewScreen()
	if error != nil {
		fmt.Fprintf(os.Stderr, "%v\n", error)
		os.Exit(ERROR_SCREEN_INIT)
	}
	if error = screen.Init(); error != nil {
		fmt.Fprintf(os.Stderr, "%v\n", error)
		os.Exit(ERROR_SCREEN_INIT)
	}
	// A digital frontier...
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		screen.Clear()
	}
}

// ExitGame is a callback triggered when the game terminates
func (t TextGame) ExitGame(playAgain bool) {
}
