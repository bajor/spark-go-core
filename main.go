package main

import (
	"fmt"

	lazy "github.com/bajor/spark-go-core/lazy_evaluation"
)

func main() {
	chain := &lazy.LazyChain{}

	chain.Add(lazy.MapOperation(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	}))

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
// write function for string only - do some manipulations on the stirng
