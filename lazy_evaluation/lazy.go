// https://chat.openai.com/c/1b4da79b-98fe-456e-8e1e-c6cdde63f161

package lazy

import (
	"fmt"
)

// Iterator interface for lazy evaluation
type Iterator interface {
	Next() (interface{}, bool)
	Reset()
}

// Base iterator that wraps a slice
type SliceIterator struct {
	data  []interface{}
	index int
}

func NewSliceIterator(data []interface{}) *SliceIterator {
	return &SliceIterator{
		data:  data,
		index: -1,
	}
}

func (si *SliceIterator) Next() (interface{}, bool) {
	si.index++
	if si.index >= len(si.data) {
		return nil, false
	}
	return si.data[si.index], true
}

func (si *SliceIterator) Reset() {
	si.index = -1
}

// Filter iterator that applies filter operations lazily
type FilterIterator struct {
	iterator Iterator
	filters  []func(interface{}) bool
}

func NewFilterIterator(iterator Iterator, filters []func(interface{}) bool) *FilterIterator {
	return &FilterIterator{
		iterator: iterator,
		filters:  filters,
	}
}

func (fi *FilterIterator) Next() (interface{}, bool) {
	for {
		item, hasNext := fi.iterator.Next()
		if !hasNext {
			return nil, false
		}

		// Apply all filters
		passes := true
		for _, filter := range fi.filters {
			if !filter(item) {
				passes = false
				break
			}
		}

		if passes {
			return item, true
		}
		// If it doesn't pass, continue to next item
	}
}

func (fi *FilterIterator) Reset() {
	fi.iterator.Reset()
}

// Map iterator that applies map operations lazily
type MapIterator struct {
	iterator Iterator
	mappers  []func(interface{}) (interface{}, error)
}

func NewMapIterator(iterator Iterator, mappers []func(interface{}) (interface{}, error)) *MapIterator {
	return &MapIterator{
		iterator: iterator,
		mappers:  mappers,
	}
}

func (mi *MapIterator) Next() (interface{}, bool) {
	item, hasNext := mi.iterator.Next()
	if !hasNext {
		return nil, false
	}

	// Apply all map operations
	var err error
	for _, mapper := range mi.mappers {
		item, err = mapper(item)
		if err != nil {
			panic(fmt.Sprintf("Map operation failed: %v", err))
		}
	}

	return item, true
}

func (mi *MapIterator) Reset() {
	mi.iterator.Reset()
}

// LazyChain represents a chain of lazy operations
type LazyChain struct {
	iterator Iterator
	filters  []func(interface{}) bool
	mappers  []func(interface{}) (interface{}, error)
	reducers []func([]interface{}) ([]interface{}, error)
}

func NewLazyChain(data []interface{}) *LazyChain {
	return &LazyChain{
		iterator: NewSliceIterator(data),
		filters:  make([]func(interface{}) bool, 0),
		mappers:  make([]func(interface{}) (interface{}, error), 0),
		reducers: make([]func([]interface{}) ([]interface{}, error), 0),
	}
}

func (lc *LazyChain) AddFilter(filter func(interface{}) bool) *LazyChain {
	lc.filters = append(lc.filters, filter)
	return lc
}

func (lc *LazyChain) AddMap(mapper func(interface{}) (interface{}, error)) *LazyChain {
	lc.mappers = append(lc.mappers, mapper)
	return lc
}

func (lc *LazyChain) AddReduce(reducer func([]interface{}) ([]interface{}, error)) *LazyChain {
	lc.reducers = append(lc.reducers, reducer)
	return lc
}

func (lc *LazyChain) HasOperations() bool {
	return len(lc.filters) > 0 || len(lc.mappers) > 0 || len(lc.reducers) > 0
}

// Collect evaluates the lazy chain and returns all results
func (lc *LazyChain) Collect() ([]interface{}, error) {
	// Build the iterator chain
	currentIterator := lc.iterator

	// Apply filters first
	if len(lc.filters) > 0 {
		currentIterator = NewFilterIterator(currentIterator, lc.filters)
	}

	// Apply maps
	if len(lc.mappers) > 0 {
		currentIterator = NewMapIterator(currentIterator, lc.mappers)
	}

	// Collect all results
	var results []interface{}
	for {
		item, hasNext := currentIterator.Next()
		if !hasNext {
			break
		}
		results = append(results, item)
	}

	// Apply reducers if any
	for _, reducer := range lc.reducers {
		var err error
		results, err = reducer(results)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ForEach applies a function to each element without collecting results
func (lc *LazyChain) ForEach(fn func(interface{}) error) error {
	currentIterator := lc.iterator

	// Apply filters first
	if len(lc.filters) > 0 {
		currentIterator = NewFilterIterator(currentIterator, lc.filters)
	}

	// Apply maps
	if len(lc.mappers) > 0 {
		currentIterator = NewMapIterator(currentIterator, lc.mappers)
	}

	// Process each element
	for {
		item, hasNext := currentIterator.Next()
		if !hasNext {
			break
		}
		if err := fn(item); err != nil {
			return err
		}
	}

	return nil
}

// Legacy compatibility methods
type Operation interface{}

type FilterOperations func(i interface{}) bool
type MapOperation func(i interface{}) (interface{}, error)
type ReduceOperation func(a []interface{}) ([]interface{}, error)

func (lc *LazyChain) Add(op Operation) {
	switch opTyped := op.(type) {
	case FilterOperations:
		lc.AddFilter(opTyped)
	case MapOperation:
		lc.AddMap(opTyped)
	case ReduceOperation:
		lc.AddReduce(opTyped)
	default:
		fmt.Println("Unsupported operation type")
	}
}

func (lc *LazyChain) Evaluate(inputs []interface{}) (interface{}, error) {
	// Reset the chain with new data
	lc.iterator = NewSliceIterator(inputs)
	return lc.Collect()
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
