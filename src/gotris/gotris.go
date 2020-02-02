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
	fmt.Printf("Hello, world!\n")
	fmt.Printf("Test %d", model.Test())
}
