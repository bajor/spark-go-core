package rdd

import (
	lazy "github.com/bajor/spark-go-core/lazy_evaluation"
)

type KeyedRDD struct {
	data  []interface{}
	chain *lazy.LazyChain
	key   func(i interface{}) (interface{}, error)
}

func NewKeyedRDD(data []interface{}, key func(i interface{}) (interface{}, error)) *KeyedRDD {
	return &KeyedRDD{
		data:  data,
		chain: &lazy.LazyChain{},
		key:   key,
	}
}

func (r *KeyedRDD) Map(f func(i interface{}) (interface{}, error)) *KeyedRDD {
	r.chain.Add(lazy.MapOperation(f))
	return r
}

func (r *KeyedRDD) Filter(f func(i interface{}) bool) *KeyedRDD {
	r.chain.Add(lazy.FilterOperations(f))
	return r
}

func (r *KeyedRDD) ReduceByKey(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	r.chain.Add(lazy.ReduceOperation(f))
	return r
}

func (r *KeyedRDD) Reduce(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	r.chain.Add(lazy.ReduceOperation(f))
	return r
}

func (r *KeyedRDD) GetData() []interface{} {
	// Evaluate the lazy chain if there are operations
	if r.chain.HasOperations() {
		result, err := r.chain.Evaluate(r.data)
		if err != nil {
			panic(err)
		}
		if resultSlice, ok := result.([]interface{}); ok {
			return resultSlice
		}
		panic("Unexpected result type from lazy evaluation")
	}
	return r.data
}
