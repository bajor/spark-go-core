package rdd

import (
	"github.com/bajor/spark-go-core/types"
)

// KeyedRDD embeds the types.KeyedRDD to allow method definitions
type KeyedRDD struct {
	*types.KeyedRDD
}

// NewKeyedRDD creates a new KeyedRDD with the given data and key function
func NewKeyedRDD(data []interface{}, key func(i interface{}) (interface{}, error)) *KeyedRDD {
	return &KeyedRDD{
		KeyedRDD: &types.KeyedRDD{
			Data:  data,
			Chain: &types.OperationChain{Operations: make([]types.Operation, 0)},
			Key:   key,
		},
	}
}

// Map applies a transformation function to each element
func (r *KeyedRDD) Map(f func(i interface{}) (interface{}, error)) *KeyedRDD {
	newChain := &types.OperationChain{Operations: make([]types.Operation, len(r.Chain.Operations))}
	copy(newChain.Operations, r.Chain.Operations)
	newChain.Operations = append(newChain.Operations, MapOperation{f: f})
	
	return &KeyedRDD{
		KeyedRDD: &types.KeyedRDD{
			Data:  r.Data,
			Chain: newChain,
			Key:   r.Key,
		},
	}
}

// Filter keeps only elements that match the predicate
func (r *KeyedRDD) Filter(f func(i interface{}) bool) *KeyedRDD {
	newChain := &types.OperationChain{Operations: make([]types.Operation, len(r.Chain.Operations))}
	copy(newChain.Operations, r.Chain.Operations)
	newChain.Operations = append(newChain.Operations, FilterOperation{f: f})
	
	return &KeyedRDD{
		KeyedRDD: &types.KeyedRDD{
			Data:  r.Data,
			Chain: newChain,
			Key:   r.Key,
		},
	}
}

// ReduceByKey groups elements by key and applies a reduce function to each group
func (r *KeyedRDD) ReduceByKey(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newChain := &types.OperationChain{Operations: make([]types.Operation, len(r.Chain.Operations))}
	copy(newChain.Operations, r.Chain.Operations)
	newChain.Operations = append(newChain.Operations, ReduceByKeyOperation{
		keyFunc:    r.Key,
		reduceFunc: f,
	})
	
	return &KeyedRDD{
		KeyedRDD: &types.KeyedRDD{
			Data:  r.Data,
			Chain: newChain,
			Key:   r.Key,
		},
	}
}

// Reduce applies a function to combine all elements into a single result
func (r *KeyedRDD) Reduce(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newChain := &types.OperationChain{Operations: make([]types.Operation, len(r.Chain.Operations))}
	copy(newChain.Operations, r.Chain.Operations)
	newChain.Operations = append(newChain.Operations, ReduceOperation{f: f})
	
	return &KeyedRDD{
		KeyedRDD: &types.KeyedRDD{
			Data:  r.Data,
			Chain: newChain,
			Key:   r.Key,
		},
	}
}

// GetData evaluates the lazy operation chain and returns the result
func (r *KeyedRDD) GetData() []interface{} {
	if len(r.Chain.Operations) == 0 {
		return r.Data
	}

	currentData := r.Data
	for _, op := range r.Chain.Operations {
		result, err := op.Execute(currentData)
		if err != nil {
			panic(err)
		}
		currentData = result
	}
	return currentData
} 