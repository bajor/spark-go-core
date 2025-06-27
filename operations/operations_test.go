package operations

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	data := []interface{}{1, 2, 3, 4}
	
	result, err := Map(data, func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	})
	
	if err != nil {
		t.Errorf("Map failed with error: %v", err)
	}
	
	expected := []interface{}{2, 4, 6, 8}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map failed: got %v, want %v", result, expected)
	}
}

func TestFilter(t *testing.T) {
	data := []interface{}{1, 2, 3, 4, 5, 6}
	
	result := Filter(data, func(i interface{}) bool {
		return i.(int) > 3
	})
	
	expected := []interface{}{4, 5, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter failed: got %v, want %v", result, expected)
	}
}

func TestReduce(t *testing.T) {
	data := []interface{}{1, 2, 3, 4}
	
	result, err := Reduce(data, func(a []interface{}) ([]interface{}, error) {
		sum := 0
		for _, val := range a {
			sum += val.(int)
		}
		return []interface{}{sum}, nil
	})
	
	if err != nil {
		t.Errorf("Reduce failed with error: %v", err)
	}
	
	expected := []interface{}{10}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Reduce failed: got %v, want %v", result, expected)
	}
}

func TestReduceByKey(t *testing.T) {
	data := []interface{}{1, 2, 3, 4, 5, 6}
	
	result, err := ReduceByKey(
		data,
		func(i interface{}) (interface{}, error) {
			return i.(int) % 2, nil // Group by even/odd
		},
		func(a []interface{}) ([]interface{}, error) {
			return []interface{}{len(a)}, nil
		},
	)
	
	if err != nil {
		t.Errorf("ReduceByKey failed with error: %v", err)
	}
	
	// Should have 2 groups: even (3 elements) and odd (3 elements)
	expected := []interface{}{3, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ReduceByKey failed: got %v, want %v", result, expected)
	}
} 