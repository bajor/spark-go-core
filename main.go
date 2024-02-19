package main

import (
	"fmt"

	lazy "github.com/bajor/spark-go-core/lazy_evaluation"
)

func main() {
	// Example usage

	chain := &lazy.LazyChain{}

	// Add operations to the chain
	chain.Add(func(i interface{}) (interface{}, error) {
		// Example operation: add 5
		return i.(int) + 5, nil
	})
	chain.Add(func(i interface{}) (interface{}, error) {
		// Another operation: multiply by 2
		return i.(int) * 2, nil
	})

	// Evaluate the chain
	input_slice := []interface{}{1, 2, 3, 4}

	result, err := chain.Evaluate(input_slice)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result)
}

// TODO:
// add fuctions like count, filter, map and introduce lazy evaluatin to them