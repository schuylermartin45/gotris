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
	"strings"
)

/***** Constants *****/

// Various gameplay modes
const (
	DEBUG_MODE string = "debug"
	TEXT_MODE  string = "text"
)

// USAGE message to display on bad input
const USAGE string = "Usage: gotris [render mode] [help]"

/***** Functions *****/

/*
 Main entry point of the Gotris project.
*/
func main() {
	// Set a default mode and construct a look-up table
	mode := TEXT_MODE
	modeMap := map[string]view.Display{
		DEBUG_MODE: new(view.DebugGame),
		TEXT_MODE:  new(view.TextGame),
	}

	// Handle user input
	argc := len(os.Args)
	if argc > 1 {
		if _, ok := modeMap[os.Args[1]]; ok {
			mode = os.Args[1]
		} else if strings.ToLower(os.Args[1]) == "help" {
			fmt.Println("Gotris: A Go-implementation of Tetris")
			fmt.Println("\nAbout")
			fmt.Println("  Author: Schuyler Martin")
			fmt.Println("  Date:   January 2020")
			fmt.Println("\n" + USAGE + "\n")
			fmt.Println("Render modes:")
			fmt.Println("  * `debug`: Basic rendering mode, used for debugging.")
			fmt.Println("  * `text`: Advanced text rendering mode.")
			os.Exit(view.EXIT_SUCCESS)
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", USAGE)
			os.Exit(view.ERROR_USAGE)
		}
		if argc > 2 {
			if strings.ToLower(os.Args[2]) == "help" {
				fmt.Println(modeMap[mode].RenderHelpMenu())
				os.Exit(view.EXIT_SUCCESS)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", USAGE)
				os.Exit(view.ERROR_USAGE)
			}
		}
	}

	// Initialize, run, and exit with the selected mode
	playAgain := true
	for playAgain {
		modeMap[mode].InitGame(model.NewBoard())
		playAgain = modeMap[mode].RenderGame()
	}
	modeMap[mode].ExitGame()
}
