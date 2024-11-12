// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Main function
// ---------------------------------------------------------------------------------------------------------------------

package main

import (
	"EmbyExplorer_for_Windows/ui"
)

func main() {
	err := ui.CreateUi()
	if err != nil {
		panic(err)
	}
}
