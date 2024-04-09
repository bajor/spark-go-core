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
	// Partition data by key
	partitionedData := make(map[interface{}][]interface{})
	for _, item := range r.data {
		key, err := r.key(item)
		if err != nil {
			panic(err)
		}
		partitionedData[key] = append(partitionedData[key], item)
	}

	// Apply reduce operation to each partition
	for key, partition := range partitionedData {
		result, err := f(partition)
		if err != nil {
			panic(err)
		}
		partitionedData[key] = result
	}

	// Flatten partitioned data back into slice
	flattenedData := make([]interface{}, 0)
	for _, partition := range partitionedData {
		flattenedData = append(flattenedData, partition...)
	}

	// Create new RDD with reduced data
	return NewKeyedRDD(flattenedData, r.key)
}

func (r *KeyedRDD) Reduce(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	// Convert RDD to slice
	dataSlice := r.GetData()

	// Apply reduce operation
	result, err := f(dataSlice)
	if err != nil {
		panic(err)
	}

	// Create new RDD with reduced data
	return NewKeyedRDD(result, r.key)
}

func (r *KeyedRDD) GetData() []interface{} {
	return r.data
}
