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
	"./view"
	"fmt"
	"os"
)

/***** Constants *****/

// Various gameplay modes
const (
	DEBUG_MODE string = "debug"
	TEXT_MODE  string = "text"
)

// USAGE message to display on bad input
const USAGE string = "Usage: gotris [render mode]"

/***** Functions *****/

/*
 Main entry point of the Gotris project.
*/
func main() {
	// Set a default mode and construct a look-up table
	mode := DEBUG_MODE
	modeMap := map[string]view.Display{
		DEBUG_MODE: new(view.DebugGame),
		TEXT_MODE:  new(view.TextGame),
	}

	// Handle user input
	if len(os.Args) > 1 {
		if _, ok := modeMap[os.Args[1]]; ok {
			mode = os.Args[1]
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", USAGE)
			os.Exit(view.ERROR_USAGE)
		}
	}

	// Initialize, run, and exit with the selected mode
	modeMap[mode].InitGame(model.NewBoard())
	modeMap[mode].RenderGame()
	// TODO implement "play again" option
	modeMap[mode].ExitGame(false)

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
