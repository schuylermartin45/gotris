/*
 * File:        debugGame.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: The most basic gameplay mode, requiring no additional
 *              dependencies. Movement is slow and the display looks awful.
 */
package view

import (
	"../model"
	"bufio"
	"fmt"
	"os"
	"strings"
)

/***** Types *****/

// KeyMap Maps keyboard input to actions.
type KeyMap map[string]Action

// DebugGame renders Gotris in a limited text environment
type DebugGame struct {
	// Model of the game
	board *model.Board
	// Input buffer
	reader *bufio.Reader
}

/***** Functions *****/

/**
Retrieves an action from a string.

@param action	String to interpret as an action

@return Action derived from the string provided. If not found, `ActionIllegal`
is returned.
*/
func getAction(action string) Action {
	action = strings.ToLower(strings.TrimSuffix(action, "\n"))
	var keyMap KeyMap = map[string]Action{
		"a":      ActionLeft,
		"left":   ActionLeft,
		"d":      ActionRight,
		"right":  ActionRight,
		"s":      ActionDown,
		"down":   ActionDown,
		"w":      ActionRotate,
		"rotate": ActionRotate,
		" ":      ActionFastDown,
		"e":      ActionExit,
		"exit":   ActionExit,
	}
	if value, ok := keyMap[action]; ok {
		return value
	}
	return ActionIllegal
}

/***** Methods *****/

// RenderHelpMenu returns a string to display the help menu in the terminal.
func (d DebugGame) RenderHelpMenu() string {
	return "Debug Mode\n" +
		"\nAbout\n" +
		"  This mode is a basic text-mode written for debugging the game.\n" +
		"  It is written only using standard Go packages.\n" +
		"\nControls\n" +
		"  * W:       Rotate\n" +
		"  * A:       Move left\n" +
		"  * S:       Move right\n" +
		"  * D:       Move down\n" +
		"  * [Space]: Drop tile to floor\n" +
		"  * E:       Exit game\n"
}

// InitGame initializes the game.
func (d *DebugGame) InitGame(b *model.Board) {
	d.board = b
	d.reader = bufio.NewReader(os.Stdin)
}

// RenderGame runs the primary gameplay loop.
func (d *DebugGame) RenderGame() {
	for {
		// Advance the game
		_, endGame := d.board.Next()

		// Draw the board
		fmt.Printf("Score:  %8v\n", d.board.GetDisplayScore())
		fmt.Println("----------------")
		d.drawItem()

		// Handle user input
		fmt.Print("Next move (w/a/s/d/ /e): ")
		keypress, _ := d.reader.ReadString('\n')
		ActionHandler(d.board, getAction(keypress), func() {
			endGame = true
		})

		// Stop the loop on the event that the game has ended.
		if endGame {
			break
		}
	}
}

// ExitGame is a callback triggered when the game terminates
func (d *DebugGame) ExitGame(playAgain bool) {
}

/** Internal **/

/*
 Dumps a tile or board to a string for printing.

 @return Dumps the game board as a simple string of 0s and 1s.
*/
func (d DebugGame) drawItem() {
	view := ""
	d.board.RenderBoard(func(row uint8, col uint8, isEOL bool, color model.TileColor) {
		// The original Tetris used 2 text characters to represent 1 unit of
		// width. After rendering each bit as 1 text character, this made a lot
		// of sense, as the the width and height now visually closer to a 1:1
		// proportion (as opposed to being closer to 1:2).
		if color == model.Transparent {
			view += "00"
		} else {
			view += "11"
		}
		// Add a newline after the last character in the row
		if isEOL {
			view += "\n"
		}
	})
	fmt.Print(view)
}
