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

// BoardGrid is one unit longer than it's displayable form. This makes collision
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

/***** Internal Functions *****/

/*
 Check collisions given a future version of the board and tile.

 @param grid		Working copy of the grid.
 @param tile		Working copy of the tile.
 @param tileDepth	Working copy of the tile depth.

 @return True if a collision was detected. False otherwise.
*/
func checkCollisions(grid BoardGrid, tile Tile, tileDepth uint8) bool {
	bottomGap := tile.GetBottomGap()
	// Take the gap at the bottom of the tile into account only if we won't
	// underflow index.
	if tileDepth > bottomGap {
		tileDepth -= bottomGap
	}
	// Detect collisions starting at the first occupied row at the bottom of the
	// tile's structure.
	bottomTileDiff := int(bottomGap) + 1
	for row := len(tile.shape) - bottomTileDiff; row >= 0; row-- {
		// If tile intersects with part of the board, a collision occurred.
		if (grid[tileDepth] & tile.shape[row]) != 0 {
			return true
		}
		// Break early to stay in bounds when part of the tile is still above the
		// screen.
		if tileDepth == 0 {
			break
		}
		tileDepth--
	}
	return false
}

/***** Methods *****/

/*
 Get the displayable version of the score.

 @return The game's current score as a displayable string
*/
func (b Board) GetDisplayScore() string {
	return fmt.Sprintf("%05d", b.score) + "00"
}

/*
 Moves the current tile to the left, if possible.

 @return True if the move happened. False otherwise.
*/
func (b *Board) MoveLeft() bool {
	return b.moveX(Left)
}

/*
 Moves the current tile to the right, if possible.

 @return True if the move happened. False otherwise.
*/
func (b *Board) MoveRight() bool {
	return b.moveX(Right)
}

/*
 Rotates the current tile, if possible.

 @return True if the move happened. False otherwise.
*/
func (b *Board) Rotate() bool {
	if b.tile == nil {
		return false
	}
	tempTile := *b.tile
	// If a tile is close to either edge, shift in the opposite direction
	// and then rotate.
	const leftBoundMask uint8 = 0b11000000
	const rightBoundMask uint8 = 0b00000011
	for row := 0; row < len(tempTile.shape); row++ {
		if (tempTile.shape[row] & leftBoundMask) > 0 {
			tempTile.MoveX(Right)
			tempTile.MoveX(Right)
			break
		} else if (tempTile.shape[row] & rightBoundMask) > 0 {
			tempTile.MoveX(Left)
			tempTile.MoveX(Left)
			break
		}
	}
	tempTile.Rotate()
	if checkCollisions(b.grid, tempTile, b.tileDepth) {
		return false
	}
	*b.tile = tempTile
	return true
}

/*
 Handle the next iteration of the game. Coupled with the primary game loop,
 this makes the game work.

 @return The current grid to display AND true if the game has ended.
*/
func (b *Board) Next() ([]uint8, bool) {
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
		return b.grid[:BoardHeight], false
	}

	// Track conditions for moving to the next tile. In other words, a collision
	// has been detected.
	tileDone := false
	// Track if the game is done ("We're in the end game now, Stark")
	gameDone := false

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
	if checkCollisions(b.grid, *b.tile, b.tileDepth+1) {
		tileDone = true
		// The game ends when a collision is detected on a tile that has yet
		// to drop into the board.
		if b.tileDepth < uint8(len(b.tile.shape)) {
			gameDone = true
		}
	}

	// Advance to the next tile. Tile becomes persistently part of the board
	if tileDone {
		b.tile = nil
		// Search for filled rows, clear them, shift above rows down.
		// Remember that there is a phantom row at the bottom of the board that is
		// not rendered.
		for row := int8(BoardHeight - 1); row >= 0; row-- {
			if workingGrid[row] == 255 {
				for i := row; i >= 1; i-- {
					workingGrid[i] = workingGrid[i-1]
				}
				// Top row gets wiped clean
				workingGrid[0] = 0
				// Score gets incremented
				b.score++
			}
		}
		b.grid = workingGrid
	} else {
		b.tileDepth++
	}
	return workingGrid[:BoardHeight], gameDone
}

/***** Internal Methods *****/

/*
 Helper function that moves in either X direction.

 @return True if the move happened. False otherwise.
*/
func (b *Board) moveX(direction XDirection) bool {
	if b.tile == nil {
		return false
	}
	tempTile := *b.tile
	tempTile.MoveX(direction)
	if checkCollisions(b.grid, tempTile, b.tileDepth) {
		return false
	}
	*b.tile = tempTile
	return true
}
