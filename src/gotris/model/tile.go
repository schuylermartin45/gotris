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

/***** Methods *****/

/*
 Picks a tile at random
*/
func (b Tile) PickTile() Tile {
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
	return tiles[rand.Intn(len(tiles))]
}

/*
 Get the color of the tile

 @return The tile's color enumeration.
*/
func (t Tile) GetTileColor() TileColor {
	return t.color
}
