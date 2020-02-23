/*
 * File:        tile.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Representation of a tile in Gotris.
 */
package model

import (
	"math/rand"
)

/***** Types *****/

// XDirection describes movement on the x-axis
type XDirection bool

// XDirection enumerations
const (
	Left  XDirection = true
	Right XDirection = false
)

// TileColor represents a color in an enumerated form
type TileColor uint8

// TileColor enumerations. Also doubles as a way to identify each tile shape.
const (
	// Also known as the "nil" color
	Transparent TileColor = 0
	Blue        TileColor = 1
	Cyan        TileColor = 2
	Grey        TileColor = 3
	Yellow      TileColor = 4
	Green       TileColor = 5
	Violet      TileColor = 6
	Red         TileColor = 7
)

// TileSize is the max width/height/number of blocks in a tile
const TileSize = 4

// SimpleBlock is the old format used to generate shapes.
type SimpleBlock [TileSize]uint8

// Block is the primitive structure that describes the shape of each tile.
type Block [TileSize]uint32

// Tile represents a tile in the game.
type Tile struct {
	// Each piece uses an array of numbers.
	shape Block
	// Color information associated with the block.
	color TileColor
}

/***** Functions *****/

/*
 Picks a tile at random

 @param random Reference to a random number generator object.
*/
func PickTile(random *rand.Rand) *Tile {
	// Tiles follow the Windows 98 Tetris Color scheme.
	tiles := [7]Tile{
		// L-left _|
		buildTile(SimpleBlock{
			0b00000000,
			0b00001000,
			0b00001000,
			0b00011000,
		}, Violet),
		// L-right |_
		buildTile(SimpleBlock{
			0b00000000,
			0b00010000,
			0b00010000,
			0b00011000,
		}, Yellow),
		// Square
		buildTile(SimpleBlock{
			0b00000000,
			0b00011000,
			0b00011000,
			0b00000000,
		}, Cyan),
		// Pipe
		buildTile(SimpleBlock{
			0b00001000,
			0b00001000,
			0b00001000,
			0b00001000,
		}, Red),
		// Tri-point _-_
		buildTile(SimpleBlock{
			0b00000000,
			0b00010000,
			0b00111000,
			0b00000000,
		}, Grey),
		// S
		buildTile(SimpleBlock{
			0b00000000,
			0b00011000,
			0b00110000,
			0b00000000,
		}, Blue),
		// Z
		buildTile(SimpleBlock{
			0b00000000,
			0b00110000,
			0b00011000,
			0b00000000,
		}, Green),
	}
	// Note to self: this is legit in Go even if it feels so wrong.
	return &tiles[random.Intn(len(tiles))]
}

/*
 Converts from the old 8-bit based grid system to the 32-bit color one (so the
 tiles can be visibly "drawn" in their binary form.

 @param shape Old shape, 8-bit representation
 @param color 3-bit color code
*/
func buildTile(shape SimpleBlock, color TileColor) Tile {
	newShape := Block{}
	for row := 0; row < TileSize; row++ {
		if shape[row] != 0 {
			var mask uint8 = 1 << 7
			for col := uint8(0); col < 8; col++ {
				tempRow := uint32(0)
				if (shape[row] & mask) > 0 {
					// Initialize with the color, set on the left-hand side
					// of the board, minding the spare right-most bit.
					//
					// << 28 gets leading bit to first position, minus left
					// 1 bit pad, - blockBitSize to include the new left-most
					// column that didn't exist in the original 8-column version
					tempRow = uint32(color) << (rShiftBlockBitDiff - blockBitSize)
					// Project the color in the new board dimensions.
					tempRow >>= col * uint8(blockBitSize)
				}
				// Accumulate blocks we've seen in the horizontal axis.
				newShape[row] |= tempRow
				mask >>= 1
			}
		}
	}
	tile := Tile{
		shape: newShape,
		color: color,
	}
	return tile
}

/***** Methods *****/

/*
 Move the tile one unit in the x-axis (left or right )
*/
func (t *Tile) MoveX(direction XDirection) {
	// Check the bounds. If the left-most or right-most bit is set in any column,
	// then we can no longer move in that direction.
	const (
		leftBoundMask  uint32 = 0xF0000000 // 1 leading unused bit + 1 block
		rightBoundMask uint32 = 0x0000000F // 1 trailing unused bit + 1 block
	)
	for row := 0; row < len(t.shape); row++ {
		if (direction == Left) && (t.shape[row]&leftBoundMask) > 0 {
			return
		} else if (direction == Right) && (t.shape[row]&rightBoundMask) > 0 {
			return
		}
	}

	// If it is safe to shift, shift tile
	for row := 0; row < len(t.shape); row++ {
		switch direction {
		case Left:
			t.shape[row] <<= blockBitSize
		case Right:
			t.shape[row] >>= blockBitSize
		}
	}
}

/*
 Rotates the tile by 90 degrees.
*/
func (t *Tile) Rotate() {
	// Short-circuit on the square tile. No rotation is required.
	if t.color == Cyan {
		return
	}
	// Generate a repeating color mask to make it easier to copy the color
	// into the transposed matrix.
	colorMask := uint32(0)
	for col := uint32(0); col < 10; col++ {
		colorMask |= uint32(t.color) << ((blockBitSize * col) + 1)
	}

	// Find all of the positions in the board that are currently filled. Track
	// the minimum column value as that minimum column becomes the first row.
	var rowIdxs []uint8
	var colIdxs []uint8
	minCol := BoardWidth
	for row := uint8(0); row < TileSize; row++ {
		var mask uint32 = blockMask << rShiftBlockBitDiff
		for col := uint8(0); col < BoardWidth; col++ {
			if uint32(t.shape[row]&mask) > 0 {
				rowIdxs = append(rowIdxs, row)
				colIdxs = append(colIdxs, col)
				if col < minCol {
					minCol = col
				}
			}
			mask >>= blockBitSize
		}
	}
	// Iterate over all known block positions, re-adjusting the coordinates
	// as blocks are examined.
	transpose := Block{}
	for i := 0; i < TileSize; i++ {
		transposeMask := uint32(blockMask << rShiftBlockBitDiff)
		transposeMask >>= blockBitSize * uint32(rowIdxs[i]+minCol)
		transpose[colIdxs[i]-minCol] |= transposeMask & colorMask
	}
	t.shape = transpose
}

/*
 Get the color of the tile

 @return The tile's color enumeration.
*/
func (t Tile) GetColor() TileColor {
	return t.color
}

/*
 Get a tile's block structure/shape

 @return The tile's shape.
*/
func (t Tile) GetBlock() []uint32 {
	return t.shape[:]
}

/*
 Get the size of the gap from the bottom of the physical tile to the end
 of the tile's block. In other words, this is the count of zero-rows at the end
 of the tile.

 @return Number of empty rows under the tile in the tile's block structure.
*/
func (t Tile) GetBottomGap() uint8 {
	var cntr uint8 = 0
	for row := len(t.shape) - 1; row >= 0; row-- {
		if t.shape[row] == 0 {
			cntr++
		} else {
			return cntr
		}
	}
	return cntr
}
