/*
 * File:        board.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Representation of the Gotris board.
 */
package model

/***** Types *****/
type Board struct {
	// Board will be 8 units wide, 20 tall. The width allows me to do some
	// fancy bitwise operations later because they're fun.
	Grid [20]uint8
}

/***** Methods *****/

/*
 Dumps a board to a string for printing.
 TODO: This should be moved into a view/rendering engine.
*/
func (b Board) DumpBoard() string {
	view := ""
	for row := 0; row < len(b.Grid); row++ {
		var mask uint8 = (1 << 7)
		for col := 0; col < 8; col++ {
			mask >>= 1
			// The original Tetris used 2 text characters to represent 1 unit of
			// width. After rendering each bit as 1 text character, this made a lot
			// of sense, as the the width and height now visually closer to a 1:1
			// proportion (as opposed to being closer to 1:2).
			if (b.Grid[row] & mask) == 1 {
				view += "11"
			} else {
				view += "00"
			}
		}
		view += "\n"
	}
	return view
}