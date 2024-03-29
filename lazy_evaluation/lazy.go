// https://chat.openai.com/c/1b4da79b-98fe-456e-8e1e-c6cdde63f161

package lazy

import (
	"fmt"
)

type Operation interface{}

type FilterOperations func(i interface{}) bool

type MapOperation func(i interface{}) (interface{}, error)

type ReduceOperation func(a []interface{}) ([]interface{}, error)

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

func (lc *LazyChain) Evaluate(inputs []interface{}) (interface{}, error) {
	var err error

	for i := 0; i < len(inputs); i++ {
		for _, op := range lc.filterOps {
			if !op(inputs[i]) {
				inputs = append(inputs[:i], inputs[i+1:]...)
			}
		}
	}

	for i := 0; i < len(inputs); i++ {
		for _, op := range lc.mapOps {
			inputs[i], err = op(inputs[i])
			if err != nil {
				return nil, err
			}
		}
	}

	for _, op := range lc.reduceOps {
		inputs, err = op(inputs)
		if err != nil {
			return nil, err
		}
	}

	return inputs, nil
}

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
