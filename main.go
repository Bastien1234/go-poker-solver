package main

import (
	"fmt"
	"pokersolver/pkg/tree"
)

func main() {
	tree := tree.NewTree()
	tree.MakeRiverTree()

	fmt.Println("program finished with code 0")
}
