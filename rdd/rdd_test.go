package rdd

import (
	"reflect"
	"testing"
)

func TestRDD_FilterMapReduce(t *testing.T) {
	rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	})

	rdd = rdd.Filter(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	})

	result := rdd.GetData()
	expected := []interface{}{10, 12}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter+Map failed: got %v, want %v", result, expected)
	}

	rdd = rdd.ReduceByKey(func(a []interface{}) ([]interface{}, error) {
		return []interface{}{len(a)}, nil
	})

	result = rdd.GetData()
	expected = []interface{}{2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ReduceByKey failed: got %v, want %v", result, expected)
	}

	rdd = rdd.Reduce(func(a []interface{}) ([]interface{}, error) {
		sum := 0
		for _, val := range a {
			sum += val.(int)
		}
		return []interface{}{sum}, nil
	})

	result = rdd.GetData()
	expected = []interface{}{2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Reduce failed: got %v, want %v", result, expected)
	}
} 