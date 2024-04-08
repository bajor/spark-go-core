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
		panic(err)
	}
	fmt.Println("Result:", result)
}


/*
type RDD struct {
	data []interface{}
	chain *lazy.LazyChain
}

func NewRDD(data []interface{}) *RDD {
	return &RDD{
		data: data,
		chain: &lazy.LazyChain{},
	}
}

func (r *RDD) Map(f func(i interface{}) (interface{}, error)) *RDD {
	r.chain.Add(lazy.MapOperation(f))
	return r
}

func (r *RDD) Filter(f func(i interface{}) bool) *RDD {
	r.chain.Add(lazy.FilterOperations(f))
	return r
}

func (r *RDD) Reduce(f func(a []interface{}) ([]interface{}, error)) *RDD {
	r.chain.Add(lazy.ReduceOperation(f))
	return r
}

func (r *RDD) Count() (int, error) {
	r.chain.Add(lazy.ReduceOperation(func(a []interface{}) ([]interface{}, error) {
		return []interface{}{len(a)}, nil
	}))
	result, err := r.chain.Evaluate(r.data)
	if err != nil {
		return 0, err
	}
	return result[0].(int), nil
}

func (r *RDD) Sum() (int, error) {
	r.chain.Add(lazy.ReduceOperation(func(a []interface{}) ([]interface{}, error) {
		sum := 0
		for _, v := range a {
			sum += v.(int)
		}
		return []interface{}{sum}, nil
	}))
	result, err := r.chain.Evaluate(r.data)
	if err != nil {
		return 0, err
	}
	return result[0].(int), nil
}

func (r *RDD) Max() (int, error) {
	r.chain.Add(lazy.ReduceOperation(func(a []interface{}) ([]interface{}, error) {
		max := a[0].(int)
		for _, v := range a {
			if v.(int) > max {
				max = v.(int)
			}
		}
		return []interface{}{max}, nil
	}))
	result, err := r.chain.Evaluate(r.data)
	if err != nil {
		return 0, err
	}
	return result[0].(int), nil
}

func (r *RDD) Min() (int, error) {
	r.chain.Add(lazy.ReduceOperation(func(a []interface{}) ([]interface{}, error) {
		min := a[0].(int)
		for _, v := range a {
			if v.(int) < min {
				min = v.(int)
			}
		}
		return []interface{}{min}, nil
	}))
	result, err := r.chain.Evaluate(r.data)
	if err != nil {
		return 0, err
	}
	return result[0].(int), nil
}
*/