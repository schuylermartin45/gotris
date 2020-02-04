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

/***** Constants *****/
const (
	// Board will be 8 units wide, 20 tall. The width allows me to do some
	// fancy bitwise operations later because they're fun.
	BoardWidth  uint8 = 8
	BoardHeight uint8 = 20
)

/***** Types *****/

// The board is one unit longer than it's displayable form. This makes collision
// detection easier.
type BoardGrid [BoardHeight + 1]uint8

// Board represents the primary state of the game.
type Board struct {
	grid BoardGrid
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

/***** Functions *****/

/*
 Constructs a Gotris board.

 @return A freshly made, artisanal, Gotris board.
*/
func NewBoard() *Board {
	b := new(Board)
	// Last grid row (which is not drawn) is full of 1s for easier
	// collision detection.
	b.grid[BoardHeight] = 255
	return b
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

 @return The current grid to display.
*/
func (b *Board) Next() []uint8 {
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
		return b.grid[:BoardHeight]
	}

	// Track conditions for moving to the next tile. In other words, a collision
	// has been detected.
	tileDone := false

	// Work from the bottom of the tile piece to the top of the tile, adding it
	// into the working copy of the grid.
	workingGrid := b.grid
	boardIdx := b.tileDepth
	bottomGap := b.tile.GetBottomGap()
	// Take the gap at the bottom of the tile into account only if we won't
	// underflow index.
	if boardIdx > bottomGap {
		boardIdx -= bottomGap
	}
	bottomTileDiff := int(bottomGap) + 1
	// Only render from the physical bottom of the tile.
	for row := len(b.tile.shape) - bottomTileDiff; row >= 0; row-- {
		// Combine the tile into the board.
		workingGrid[boardIdx] |= b.tile.shape[row]
		// Break early to stay in bounds when part of the tile is still above the
		// screen.
		if boardIdx == 0 {
			break
		}
		boardIdx--
	}

	// If a collision is detected in the next move, then we stop here and move
	// to the next tile.
	if b.checkCollisions() {
		tileDone = true
	}

	// Advance to the next tile. Tile becomes persistently part of the board
	if tileDone {
		b.tile = nil
		b.grid = workingGrid
		// Advance the tile in the board
	} else {
		b.tileDepth++
	}
	return workingGrid[:BoardHeight]
}

/*
 Check collisions on the next move.

 @return A collision type. TODO make enum
*/
func (b Board) checkCollisions() bool {
	workingGrid := b.grid
	bottomGap := b.tile.GetBottomGap()
	// Advance one more unit than the game currently is at.
	var boardIdx uint8 = b.tileDepth + 1
	// Take the gap at the bottom of the tile into account only if we won't
	// underflow index.
	if boardIdx > bottomGap {
		boardIdx -= bottomGap
	}
	// Detect collisions starting at the first occupied row at the bottom of the
	// tile's structure.
	bottomTileDiff := int(bottomGap) + 1
	for row := len(b.tile.shape) - bottomTileDiff; row >= 0; row-- {
		// If tile intersects with part of the board, a collision occurred.
		if (workingGrid[boardIdx] & b.tile.shape[row]) != 0 {
			return true
		}
		// Break early to stay in bounds when part of the tile is still above the
		// screen.
		if boardIdx == 0 {
			break
		}
		boardIdx--
	}
	return false
}
