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
	// Board will be 10 units wide, 20 tall. Initially the board was 8 long
	// to performs some fun bit twiddling but I decided to expand it when I got
	// tired of the lack of colors.
	BoardWidth  uint8 = 10
	BoardHeight uint8 = 20

	/** Internal **/

	// A full row
	maskFullRow uint32 = 0xFFFFFFFF
	// An empty row (with 2 bits unused)
	maskRow2BitPad uint32 = 0x80000001
	// Bit-size of one color-block
	blockBitSize uint32 = 3
	// Amount to shift a color or mask value to the right by to be in
	// the leading (right-most) position in the board (1 bit left of the right
	// most pad)
	rShiftBlockBitDiff uint32 = 28
	// Mask used to detect blocks (will require left shifting to maintain
	// bit interpretation order)
	blockMask = 0b111
)

/***** Types *****/

// BoardGrid is one unit taller than it's displayable form. This makes collision
// detection easier.
type BoardGrid [BoardHeight + 1]uint32

/*
 DrawBlock is a callback that renders a single block when called by
 `RenderBoard()`.

 @param row   Row position in the board.
 @param col   Column position in the board.
 @param isEOL Flag indicates if this is the last column drawn in a row.
 @param color Color of the block at position (row, col).
*/
type DrawBlock func(row uint8, col uint8, isEOL bool, color TileColor)

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
	// Since we have 2 bits we can't do anything with, we pad each side
	// of the board by 1 bit
	for i := 0; i < int(BoardHeight); i++ {
		b.grid[i] = maskRow2BitPad
	}
	// Last grid row (which is not drawn) is full of 1s for easier
	// collision detection.
	b.grid[BoardHeight] = maskFullRow
	return b
}

/***** Internal Functions *****/

/*
 Helper function that calculates a "collision row", which is a row that
 represents all blocks as `0b111`, which are completely "filled in".

 Gaps in color codes can result in tiles failing to colide. For example, colors
 `0b101` and `0b010` will result in `0b000`, which will allow tiles to collide.
 Filling in all color codes as `0b111`, eliminates the issue.

 @param row Original row full of color data.

 @return The original row, but all color data replaced with "full" blocks,
         (value: `0b11`)
*/
func calcCollisionRow(row uint32) uint32 {
	var mask uint32 = blockMask << rShiftBlockBitDiff
	collisionRow := uint32(0)
	for col := uint8(0); col < BoardWidth; col++ {
		if (mask & row) > 0 {
			collisionRow |= mask
		}
		mask >>= blockBitSize
	}
	// Add the bounding bits last, so the boolean check during construction
	// of the collision row works.
	collisionRow |= maskRow2BitPad
	return collisionRow
}

/*
 Check collisions given a future version of the board and tile.

 @param grid      Working copy of the grid.
 @param tile      Working copy of the tile.
 @param tileDepth Working copy of the tile depth.

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
		collisionRow := calcCollisionRow(grid[tileDepth])
		// If tile intersects with part of the board, a collision occurred.
		//
		// This check against 0 is valid, as the bits set in `maskRow2BitPad`
		// are not set in the dropping tile.
		if (collisionRow & tile.shape[row]) != 0 {
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

/*
 Helper function that iterates over blocks, calling a render callback at each
 (row, col) position.

 @param draw	Callback to draw a block at a row, column position with a specific
             	color.
 @param blocks	Array of blocks to render.
 @param height Height of the blocks array.
 @param width  Width of the blocks array. If this is shorter than `BoardWidth`,
               the tile will attempt to be vertically centered
*/
func renderBlocks(draw DrawBlock, blocks []uint32, height uint8, width uint8) {
	// Padding calculation for width
	widthDiff := uint8(0)
	if BoardWidth > width {
		widthDiff = BoardWidth - width
	} else if BoardWidth < width {
		widthDiff = 0
	}
	halfWidthDiff := widthDiff / 2
	for row := uint8(0); row < height; row++ {
		var mask uint32 = blockMask << (rShiftBlockBitDiff - (blockBitSize * uint32(widthDiff/2)))
		paddedWidth := width + halfWidthDiff
		for col := uint8(halfWidthDiff); col < paddedWidth; col++ {
			// Select one block at a time, determine the color
			color := Transparent
			// Non-zero values require additional shifting
			singleBlock := uint32(blocks[row] & mask)
			if singleBlock > 0 {
				// Shift to the far right, so the bit can be interpretted as a
				// color. +1 is for the right-most extra bit.
				shiftBy := (blockBitSize * uint32((BoardWidth-1)-col)) + 1
				color = TileColor(singleBlock >> shiftBy)
			}
			isEOL := col >= (paddedWidth - 1)
			draw(row, col, isEOL, color)
			mask >>= blockBitSize
		}
	}
}

/***** Methods *****/

/*
 Get the displayable version of the score.

 @return The game's current score as a displayable string
*/
func (b Board) GetDisplayScore() string {
	return fmt.Sprintf("%06d", b.score) + "00"
}

/*
 Get the current level. The higher the level, the fast the game.

 For added fun (and in the spirit of Pacman) the level counter will be 8 bits
 longs. So if someone manages to get it that high, they'll start back at level
 one.

 @return The game's current level.
*/
func (b Board) GetLevel() uint8 {
	// Every ten cleared rows gets new level.
	return uint8(b.score / 10)
}

/*
 Get the next tile (for preview rendering purposes)

 @return A copy of the next tile for rendering
*/
func (b Board) GetNextTile() Tile {
	// If nil, return an empty tile
	if b.nextTile == nil {
		return Tile{}
	}
	return *b.nextTile
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
 Moves the tile down one additional unit, if possible.

 @return True if the move happened. False otherwise.
*/
func (b *Board) MoveDown() bool {
	if b.tile == nil {
		return false
	}
	tempDepth := b.tileDepth + 1
	if checkCollisions(b.grid, *b.tile, tempDepth) {
		return false
	}
	b.tileDepth = tempDepth
	return true
}

/*
 Moves the tile down until a colission occurs
*/
func (b *Board) MoveFastDown() {
	if b.tile == nil {
		return
	}
	for b.MoveDown() {
	}
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
	const (
		leftBoundMask  uint32 = 0xF7000000 // 1 leading unused bit + (2*blockBitSize) = 7 bits
		rightBoundMask uint32 = 0x0000007F // 1 trailing unused bit + (2*blockBitSize) = 7 bits
	)
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
func (b *Board) Next() ([]uint32, bool) {
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

	// Calculate the current state of the grid.
	workingGrid := b.calcWorkingGrid()
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
		numCleared := uint16(0)
		for row := int8(BoardHeight - 1); row >= 0; row-- {
			if calcCollisionRow(workingGrid[row]) == maskFullRow {
				for i := row; i >= 1; i-- {
					workingGrid[i] = workingGrid[i-1]
				}
				// Top row gets wiped clean.
				workingGrid[0] = maskRow2BitPad
				// Count the cleared rows.
				numCleared++
				// Reset row calculation to run against the same row again.
				// In the event that multiple rows are cleared at once, this
				// prevents us from leaving a full row beind.
				row++
			}
		}
		// Get a score multiplier if multiple rows are cleared at once.
		b.score += numCleared * numCleared
		b.grid = *workingGrid
	} else {
		b.tileDepth++
	}
	return workingGrid[:BoardHeight], gameDone
}

/*
 Get the current state of the board, without moving to the next iteration.

 @return The current grid to display.
*/
func (b Board) Current() []uint32 {
	// If no tile is set, then the working grid is all that is needed to be
	// displayed.
	if b.tile == nil {
		return b.grid[:BoardHeight]
	}

	return b.calcWorkingGrid()[:BoardHeight]
}

/*
 Given a callback, this function iterates over the board and executes the
 the callback to render a block on the board.

 @param draw Callback to draw a block at a row, column position with a specific
             color.
*/
func (b Board) RenderBoard(draw DrawBlock) {
	renderBlocks(draw, b.Current(), BoardHeight, BoardWidth)
}

/*
 Given a callback, this function iterates over the next tile and executes the
 the callback to render a block.

 @param draw Callback to draw a block at a row, column position with a specific
             color.
*/
func (b Board) RenderNextTile(draw DrawBlock) {
	blocks := b.GetNextTile().shape
	renderBlocks(draw, blocks[:], TileSize, BoardWidth-2)
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

/*
 Calculate the "working grid". This is the board with the current dropping
 tile merged with the remaining tile pieces. This is also the visible component
 that is accessible to the views.

 @return Working version of the grid.
*/
func (b Board) calcWorkingGrid() *BoardGrid {
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
		// Break early to stay in bounds when part of the tile is still above
		// the screen.
		if boardIdx == 0 {
			break
		}
		boardIdx--
	}
	return &workingGrid
}
