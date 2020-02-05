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
	// Ticking away, the moments that make up the dull day...
	"time"
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

// TileColor enumerations
const (
	Blue   TileColor = 0
	Cyan   TileColor = 1
	Grey   TileColor = 2
	Yellow TileColor = 3
	Green  TileColor = 4
	Violet TileColor = 5
	Red    TileColor = 6
)

// Block is the primitive structure that describes the shape of each tile.
type Block [4]uint8

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
*/
func PickTile() *Tile {
	// Tiles follow the Windows 98 Tetris Color scheme.
	tiles := [7]Tile{
		// L-left _|
		Tile{
			shape: Block{
				0b00000000,
				0b00001000,
				0b00001000,
				0b00011000,
			},
			color: Violet,
		},
		// L-right |_
		Tile{
			shape: Block{
				0b00000000,
				0b00010000,
				0b00010000,
				0b00011000,
			},
			color: Yellow,
		},
		// Square
		Tile{
			shape: Block{
				0b00000000,
				0b00011000,
				0b00011000,
				0b00000000,
			},
			color: Cyan,
		},
		// Pipe
		Tile{
			shape: Block{
				0b00001000,
				0b00001000,
				0b00001000,
				0b00001000,
			},
			color: Red,
		},
		// Tri-point _-_
		Tile{
			shape: Block{
				0b00000000,
				0b00010000,
				0b00111000,
				0b00000000,
			},
			color: Grey,
		},
		// S
		Tile{
			shape: Block{
				0b00000000,
				0b00011000,
				0b00110000,
				0b00000000,
			},
			color: Blue,
		},
		// Z
		Tile{
			shape: Block{
				0b00000000,
				0b00110000,
				0b00011000,
				0b00000000,
			},
			color: Green,
		},
	}
	ranNum := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Note to self: this is legit in Go even if it feels so wrong.
	return &tiles[ranNum.Intn(len(tiles))]
}

/***** Methods *****/

/*
 Move the tile one unit in the x-axis (left or right )
*/
func (t *Tile) MoveX(direction XDirection) {
	// Check the bounds. If the left-most or right-most bit is set in any column,
	// then we can no longer move in that direction.
	var leftBoundMask uint8 = 0b10000000
	var rightBoundMask uint8 = 0b00000001
	for row := 0; row < len(t.shape); row++ {
		if (direction == Left) && (t.shape[row]&leftBoundMask) > 0 {
			return
		} else if (direction == Right) && (t.shape[row]&rightBoundMask) > 0 {
			return
		}
	}

	// If it is safe to shift, shift tile
	for row := 0; row < len(t.shape); row++ {
		if direction == Left {
			t.shape[row] <<= 1
		} else if direction == Right {
			t.shape[row] >>= 1
		}
	}
}

/*
 Rotates the tile by 90 degrees.
*/
func (t *Tile) Rotate() {
	// Each row (byte in the original) becomes a column in the transpose.
	// This is a little tricky with the binary representation of the shapes
	// but still feasible. We focus on the inner 4x4 grid (each byte is padded
	// by two bits on the left and right sides) and calculate 2 masks that
	// shift in opposite directions.
	temp := Block{}
	var transposeMask uint8 = 0b00000100
	halfWidth := int(BoardWidth / 2)
	for row := 0; row < len(t.shape); row++ {
		var mask uint8 = 0b00100000
		for col := 0; col < halfWidth; col++ {
			if (t.shape[row] & mask) > 0 {
				temp[col] |= transposeMask
			}
			mask >>= 1
		}
		transposeMask <<= 1
	}
	t.shape = temp
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
func (t Tile) GetBlock() []uint8 {
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
