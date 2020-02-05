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
		" ":      ActionRotate,
		"e":      ActionExit,
		"exit":   ActionExit,
	}
	if value, ok := keyMap[action]; ok {
		return value
	}
	return ActionIllegal
}

/*
 Dumps a tile or board to a string for printing.

 @return Dumps the game board as a simple string of 0s and 1s.
*/
func drawItem(toDraw []uint8) {
	view := ""
	for row := 0; row < len(toDraw); row++ {
		var mask uint8 = 1 << 7
		for col := 0; col < int(model.BoardWidth); col++ {
			// The original Tetris used 2 text characters to represent 1 unit of
			// width. After rendering each bit as 1 text character, this made a lot
			// of sense, as the the width and height now visually closer to a 1:1
			// proportion (as opposed to being closer to 1:2).
			if (toDraw[row] & mask) > 0 {
				view += "11"
			} else {
				view += "00"
			}
			mask >>= 1
		}
		view += "\n"
	}
	fmt.Print(view)
}

/***** Methods *****/

// InitGame initializes the game.
func (d *DebugGame) InitGame(b *model.Board) {
	d.board = b
	d.reader = bufio.NewReader(os.Stdin)
}

// RenderGame runs the primary gameplay loop.
func (d *DebugGame) RenderGame() {
	for {
		// Advance the game
		grid, endGame := d.board.Next()

		// Draw the board
		fmt.Printf("Score:  %8v\n", d.board.GetDisplayScore())
		fmt.Println("----------------")
		drawItem(grid)

		// Handle user input
		fmt.Print("Next move (w/a/s/d/ /e): ")
		keypress, _ := d.reader.ReadString('\n')
		switch action := getAction(keypress); action {
		case ActionLeft:
			d.board.MoveLeft()
		case ActionRight:
			d.board.MoveRight()
		// TODO implement
		//case ActionDown:
		case ActionRotate:
			d.board.Rotate()
		case ActionExit:
			endGame = true
		}

		// Stop the loop on the event that the game has ended.
		if endGame {
			break
		}
	}
}

// ExitGame is a callback triggered when the game terminates
func (d *DebugGame) ExitGame(playAgain bool) {
}
