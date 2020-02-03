/*
 * File:        board.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Representation of the Gotris board.
 */
package model

import (
	"fmt"
)

/***** Types *****/

// Board represents the primary state of the game.
type Board struct {
	// Board will be 8 units wide, 20 tall. The width allows me to do some
	// fancy bitwise operations later because they're fun.
	grid [20]uint8
	// Holds the base score. Display score is this value x100 (to look cooler)
	score uint16
	// Reference to the current dropping tile. Nil means a new tile should be
	// picked.
	tile *Tile
	// Next tile is tracked for fancier displays that show the next tile as the
	// current one is dropping. Next becomes the tile dropping
	nextTile *Tile
	// Depth tracks how far down the current tile is in the board. 0 Means
	// no tile has dropped.
	tileDepth uint8
}

/***** Methods *****/

/*
 Get the displayable version of the score.

 @return The game's current score as a displayable string
*/
func (b Board) GetDisplayScore() string {
	return fmt.Sprintf("%5d", b.score) + "00"
}

/*
 Handle the next iteration of the game. Coupled with the primary game loop,
 this makes the game work.
*/
func (b *Board) Next() {
	// Initialize the next tile. This should a 1-time cost on first starting the
	// game. This simplifies the logic for setting the active tile.
	if b.nextTile == nil {
		b.nextTile = PickTile()
	}
	// On completion of a move, the next tile becomes the active and a new next
	// is picked.
	if b.tile == nil {
		b.tile = b.nextTile
		b.nextTile = PickTile()
		b.tileDepth = 0
		// Skip the rest of this iteration to give the user a break. Also ensures
		// that the `tileDepth` variable stays "in sync" with the actual row array
		// index.
		return
	}

	// Track conditions for moving to the next tile
	tileDone := false

	// TODO copy tile into the board.

	// The tile has hit rock-bottom.
	if b.tileDepth >= uint8(len(b.grid)) {
		tileDone = false
	}

	// Advance
	if tileDone {
		b.tile = nil
	} else {
		// Advance the tile in the board
		b.tileDepth++
	}
}

/*
 Dumps a board to a string for printing.
 TODO: This should be moved into a view/rendering engine.

 @return Dumps the game board as a simple string of 0s and 1s.
*/
func (b Board) DumpBoard() string {
	view := ""
	for row := 0; row < len(b.grid); row++ {
		var mask uint8 = 1
		for col := 0; col < 8; col++ {
			// The original Tetris used 2 text characters to represent 1 unit of
			// width. After rendering each bit as 1 text character, this made a lot
			// of sense, as the the width and height now visually closer to a 1:1
			// proportion (as opposed to being closer to 1:2).
			if (b.grid[row] & mask) > 0 {
				view += "11"
			} else {
				view += "00"
			}
			mask <<= 1
		}
		view += "\n"
	}
	return view
}
