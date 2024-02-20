// https://chat.openai.com/c/1b4da79b-98fe-456e-8e1e-c6cdde63f161

package lazy

import "fmt"

type Operation interface{}

type FilterOperations func(i interface{}) (interface{}, error)

type MapOperation func(i interface{}) (interface{}, error)

type ReduceOperation func(i interface{}) (interface{}, error)

// LazyChain holds a slice of operations to be executed lazily
type LazyChain struct {
	filterOps []FilterOperations
	mapOps    []MapOperation
	reduceOps []ReduceOperation
}

func (lc *LazyChain) Add(op Operation) {
	switch opTyped := op.(type) {
	case FilterOperations:
		lc.filterOps = append(lc.filterOps, opTyped)
	case MapOperation:
		lc.mapOps = append(lc.mapOps, opTyped)
	case ReduceOperation:
		lc.reduceOps = append(lc.reduceOps, opTyped)
	default:
		fmt.Println("Unsupported operation type")
	}
}

// Evaluate executes the chain of operations lazily
func (lc *LazyChain) Evaluate(inputs []interface{}) (interface{}, error) {
	var err error

	for i := 0; i < len(inputs); i++ {
		for _, op := range lc.filterOps {
			inputs[i], err = op(inputs[i])
			// if not matching the cripteria - remove it from input list
			if err != nil {
				return nil, err // Early return on first error
			}
		}
	}

	for i := 0; i < len(inputs); i++ {
		for _, op := range lc.mapOps {
			inputs[i], err = op(inputs[i])
			if err != nil {
				return nil, err // Early return on first error
			}
		}

		for _, op := range lc.reduceOps {
			inputs[i], err = op(inputs[i])
			if err != nil {
				return nil, err // Early return on first error
			}
		}

	}
	return inputs, nil
}

// func EvaluateMap

// func EvaluateReduce

// func OptimizeOperationsOrder

/*

there are two types of operation map and reduce
	if map - you apply operation to all elements individualy
	if reduce - you gather all elements and then perform operation on all elements as a whole

TODO:
1 - write generic Evaluate method which:
	a - calls OptimizeOperationsOrder() -> take a look at all operations in the stack and optimized it - return new ordered list of operations
		- optimization for now - first filters, then all maps then all reduces
2 - operations method should have embedded information wheter it's map or reduce - write common interface for tmem

*/
