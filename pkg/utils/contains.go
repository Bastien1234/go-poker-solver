package utils

import (
	"pokersolver/pkg/node"
)

func Contains(array []string, card string) bool {
	for _, el := range array {
		if el == card {
			return true
		}
	}

	return false
}

func ContainsLocalActionMap(hashmap map[string]*node.SubNode, stringHand string) bool {
	for key, _ := range hashmap {
		if key == stringHand {
			return true
		}
	}

	return false
}
