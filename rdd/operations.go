package rdd

import (
	"github.com/bajor/spark-go-core/operations"
)

// MapOperation represents a map transformation
type MapOperation struct {
	f func(interface{}) (interface{}, error)
}

func (m MapOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.Map(data, m.f)
}

// FilterOperation represents a filter transformation
type FilterOperation struct {
	f func(interface{}) bool
}

func (f FilterOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.Filter(data, f.f), nil
}

// ReduceOperation represents a reduce transformation
type ReduceOperation struct {
	f func([]interface{}) ([]interface{}, error)
}

func (r ReduceOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.Reduce(data, r.f)
}

// ReduceByKeyOperation represents a reduceByKey transformation
type ReduceByKeyOperation struct {
	keyFunc    func(interface{}) (interface{}, error)
	reduceFunc func([]interface{}) ([]interface{}, error)
}

func (r ReduceByKeyOperation) Execute(data []interface{}) ([]interface{}, error) {
	return operations.ReduceByKey(data, r.keyFunc, r.reduceFunc)
} 