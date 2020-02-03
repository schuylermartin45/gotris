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

/***** Tile Constants *****/

/***** Functions *****/

/*
 Picks a tile at random
*/
func PickTile() Tile {
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
	return tiles[ranNum.Intn(len(tiles))]
}

/***** Methods *****/

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
	for row := 0; row < len(t.shape); row++ {
		var mask uint8 = 0b00100000
		for col := 0; col < 4; col++ {
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
 Dumps a tile to a string for printing.
 TODO: This should be moved into a view/rendering engine and consolidated with
 `Board::DumpBoard()`.

 @return Dumps the game board as a simple string of 0s and 1s.
*/
func (t Tile) DumpTile() string {
	view := ""
	for row := 0; row < len(t.shape); row++ {
		var mask uint8 = 1
		for col := 0; col < 8; col++ {
			if (t.shape[row] & mask) > 0 {
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
