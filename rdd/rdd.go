package rdd

// NewKeyedRDD creates a new KeyedRDD with the given data and key function
func NewKeyedRDD(data []interface{}, key func(i interface{}) (interface{}, error)) *KeyedRDD {
	return &KeyedRDD{
		data:  data,
		chain: &OperationChain{operations: make([]Operation, 0)},
		key:   key,
	}
}

// Map applies a transformation function to each element
func (r *KeyedRDD) Map(f func(i interface{}) (interface{}, error)) *KeyedRDD {
	newChain := &OperationChain{operations: make([]Operation, len(r.chain.operations))}
	copy(newChain.operations, r.chain.operations)
	newChain.operations = append(newChain.operations, MapOperation{f: f})
	
	return &KeyedRDD{
		data:  r.data,
		chain: newChain,
		key:   r.key,
	}
}

// Filter keeps only elements that match the predicate
func (r *KeyedRDD) Filter(f func(i interface{}) bool) *KeyedRDD {
	newChain := &OperationChain{operations: make([]Operation, len(r.chain.operations))}
	copy(newChain.operations, r.chain.operations)
	newChain.operations = append(newChain.operations, FilterOperation{f: f})
	
	return &KeyedRDD{
		data:  r.data,
		chain: newChain,
		key:   r.key,
	}
}

// ReduceByKey groups elements by key and applies a reduce function to each group
func (r *KeyedRDD) ReduceByKey(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newChain := &OperationChain{operations: make([]Operation, len(r.chain.operations))}
	copy(newChain.operations, r.chain.operations)
	newChain.operations = append(newChain.operations, ReduceByKeyOperation{
		keyFunc:    r.key,
		reduceFunc: f,
	})
	
	return &KeyedRDD{
		data:  r.data,
		chain: newChain,
		key:   r.key,
	}
}

// Reduce applies a function to combine all elements into a single result
func (r *KeyedRDD) Reduce(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newChain := &OperationChain{operations: make([]Operation, len(r.chain.operations))}
	copy(newChain.operations, r.chain.operations)
	newChain.operations = append(newChain.operations, ReduceOperation{f: f})
	
	return &KeyedRDD{
		data:  r.data,
		chain: newChain,
		key:   r.key,
	}
}

// GetData evaluates the lazy operation chain and returns the result
func (r *KeyedRDD) GetData() []interface{} {
	if len(r.chain.operations) == 0 {
		return r.data
	}

	currentData := r.data
	for _, op := range r.chain.operations {
		result, err := op.Execute(currentData)
		if err != nil {
			panic(err)
		}
		currentData = result
	}
	return currentData
} 