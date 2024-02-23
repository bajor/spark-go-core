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

	chain.Add(lazy.FilterOperations(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	}))

	chain.Add(lazy.ReduceOperation(func(a []interface{}) ([]interface{}, error) {
		return []interface{}{len(a)}, nil
	}))

	input_slice := []interface{}{1, 2, 3, 4}

	result, err := chain.Evaluate(input_slice)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result)
}

// TODO:
// filter - any given filter
// map - any lambda
// reduce:
// 		- reduce - take two arguments and return one - do it cumulatively onver entire array:
//			implementations of reduce:
//			- sum
//			- dot product
//			- count - add 1 for each element of array
//			- max (iteratively max of each two), min
//
//
// add fuctions like count, filter, map and introduce lazy evaluatin to them
// write function for string only - do some manipulations on the stirng
