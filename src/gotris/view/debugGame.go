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
	"fmt"
)

/***** Types *****/

// KeyMap Maps keyboard input to actions.
type KeyMap map[string]Action

// DebugGame renders Gotris in a limited text environment
type DebugGame struct {
	board *model.Board
}

/***** Functions *****/

/**
Retrieves an action from a string.

@return Action derived from the string provided. If not found, `ActionIllegal`
is returned.
*/
func getAction(action string) Action {
	var keyMap KeyMap = map[string]Action{
		"a": ActionLeft,
		"d": ActionRight,
		"s": ActionDown,
		"w": ActionRotate,
		"e": ActionExit,
	}
	if value, ok := keyMap[action]; ok {
		return value
	}
	return 0
	//return ActionIllegal
}

/*
 Dumps a tile or board to a string for printing.

 @return Dumps the game board as a simple string of 0s and 1s.
*/
func drawItem(toDraw []uint8) {
	view := ""
	for row := 0; row < len(toDraw); row++ {
		var mask uint8 = 1
		for col := 0; col < 8; col++ {
			// The original Tetris used 2 text characters to represent 1 unit of
			// width. After rendering each bit as 1 text character, this made a lot
			// of sense, as the the width and height now visually closer to a 1:1
			// proportion (as opposed to being closer to 1:2).
			if (toDraw[row] & mask) > 0 {
				view += "11"
			} else {
				view += "00"
			}
			mask <<= 1
		}
		view += "\n"
	}
	fmt.Println(view)
}

/***** Methods *****/

// InitGame initializes the game.
func (d DebugGame) InitGame(b *model.Board) {
	d.board = b
}

// RenderGame runs the primary gameplay loop.
func (d DebugGame) RenderGame() {
	for {
	}
}

// ExitGame is a callback triggered when the game terminates
func (d DebugGame) ExitGame(playAgain bool) {
}
