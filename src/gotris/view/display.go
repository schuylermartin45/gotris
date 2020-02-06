/*
 * File:        display.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Display interface that defines how a display mode operates.
 */
package view

import (
	"../model"
)

/***** Constants *****/

// Error code constants
const (
	ERROR_USAGE       = 1
	ERROR_SCREEN_INIT = 2
)

/***** Types *****/

// Action describes a user-caused event in the game.
type Action uint8

// Enumeration of actions
const (
	ActionIllegal  Action = 0
	ActionLeft     Action = 1
	ActionRight    Action = 2
	ActionDown     Action = 3
	ActionFastDown Action = 4
	ActionRotate   Action = 5
	ActionExit     Action = 6
)

// Display is an interface that describes the features of a way to render the
// game.
type Display interface {
	// Initializes the game.
	InitGame(b *model.Board)
	// Runs the primary gameplay loop.
	RenderGame()
	// Callback for when the game terminates, with the option to play again.
	ExitGame(playAgain bool)
}
