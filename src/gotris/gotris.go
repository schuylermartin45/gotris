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
 Main entry point of the Gotris project.
*/
func main() {
	// A digital frontier...
	theGrid := model.Board{}
	fmt.Printf(theGrid.DumpBoard())
}
