package rdd

import (
	"github.com/bajor/spark-go-core/operations"
)

type KeyedRDD struct {
	data  []interface{}
	chain *OperationChain
	key   func(i interface{}) (interface{}, error)
}

type OperationChain struct {
	operations []Operation
}

type Operation interface {
	Execute(data []interface{}) ([]interface{}, error)
}

type MapOperation struct {
	f func(interface{}) (interface{}, error)
}

func (m MapOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.Map(data, m.f)
}

type FilterOperation struct {
	f func(interface{}) bool
}

func (f FilterOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.Filter(data, f.f), nil
}

type ReduceOperation struct {
	f func([]interface{}) ([]interface{}, error)
}

func (r ReduceOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.Reduce(data, r.f)
}

type ReduceByKeyOperation struct {
	keyFunc func(interface{}) (interface{}, error)
	reduceFunc func([]interface{}) ([]interface{}, error)
}

func (r ReduceByKeyOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.ReduceByKey(data, r.keyFunc, r.reduceFunc)
}

func NewKeyedRDD(data []interface{}, key func(i interface{}) (interface{}, error)) *KeyedRDD {
	return &KeyedRDD{
		data:  data,
		chain: &OperationChain{operations: make([]Operation, 0)},
		key:   key,
	}
}

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

func (r *KeyedRDD) ReduceByKey(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newChain := &OperationChain{operations: make([]Operation, len(r.chain.operations))}
	copy(newChain.operations, r.chain.operations)
	newChain.operations = append(newChain.operations, ReduceByKeyOperation{
		keyFunc: r.key,
		reduceFunc: f,
	})
	
	return &KeyedRDD{
		data:  r.data,
		chain: newChain,
		key:   r.key,
	}
}

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