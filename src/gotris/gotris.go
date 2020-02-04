/*
 * File:        gotris.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Main execution point of the `gotris` project.
 */
package main

import (
	"./model"
	"fmt"
)

/***** Functions *****/

/*
 Dumps a tile or board to a string for printing.
 TODO: This should be moved into a view/rendering engine

 @return Dumps the game board as a simple string of 0s and 1s.
*/
func drawItem(toDraw []uint8) {
	view := ""
	for row := 0; row < len(toDraw); row++ {
		var mask uint8 = 1
		for col := 0; col < 8; col++ {
			// The original Tetris used 2 text characters to represent 1 unit of
			// width. After rendering each bit as 1 text character, this made a lot
			// of sense, as the the width and height now visually closer to a 1:1
			// proportion (as opposed to being closer to 1:2).
			if (toDraw[row] & mask) > 0 {
				view += "11"
			} else {
				view += "00"
			}
			mask <<= 1
		}
		view += "\n"
	}
	fmt.Println(view)
}

/*
 Main entry point of the Gotris project.
*/
func main() {
	// A digital frontier...
	theGrid := model.NewBoard()
	for i := 0; i < 50; i++ {
		drawItem(theGrid.Next())
		fmt.Println("-------------")
	}

	// TODO implement real test cases!
	/*
		// Basic tile rotaton test
		aTile := model.PickTile()
		for i := 0; i < 5; i++ {
			drawItem(aTile.GetBlock())
			fmt.Println("-------------")
			aTile.Rotate()
		}
		// Basic Tile Movement test
		for i := 0; i < 5; i++ {
			drawItem(aTile.GetBlock())
			fmt.Println("-------------")
			aTile.MoveX(model.Left)
		}
		for i := 0; i < 10; i++ {
			drawItem(aTile.GetBlock())
			fmt.Println("-------------")
			aTile.MoveX(model.Right)
		}
	*/
}
