package rdd

import (
	"reflect"
	"testing"
)

func TestRDD_LazyEvaluation(t *testing.T) {
	// Test that operations are truly lazy - no evaluation until GetData()
	rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	// Chain multiple operations without evaluation
	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	})

	rdd = rdd.Filter(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	})

	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})

	// Now evaluate - should apply all operations in order
	result := rdd.GetData()
	expected := []interface{}{7, 9, 11, 13} // (2*3+1), (2*4+1), (2*5+1), (2*6+1)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Lazy evaluation failed: got %v, want %v", result, expected)
	}
}

func TestRDD_OperationChaining(t *testing.T) {
	rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	// Test chaining multiple operations
	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	}).Filter(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	}).Map(func(i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})

	result := rdd.GetData()
	expected := []interface{}{7, 9, 11, 13}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Operation chaining failed: got %v, want %v", result, expected)
	}
}

func TestRDD_ReduceOperations(t *testing.T) {
	rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	}).Filter(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	}).ReduceByKey(func(a []interface{}) ([]interface{}, error) {
		return []interface{}{len(a)}, nil
	})

	result := rdd.GetData()
	expected := []interface{}{1, 1, 1, 1} // Each element gets its own partition
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
	expected = []interface{}{4} // sum of [1,1,1,1]
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Reduce failed: got %v, want %v", result, expected)
	}
}

func TestRDD_ImmutableOperations(t *testing.T) {
	// Test that operations don't modify the original RDD
	original := NewKeyedRDD([]interface{}{1, 2, 3}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	mapped := original.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	})

	filtered := original.Filter(func(i interface{}) bool {
		return i.(int) > 1
	})

	// Original should be unchanged
	originalResult := original.GetData()
	expectedOriginal := []interface{}{1, 2, 3}
	if !reflect.DeepEqual(originalResult, expectedOriginal) {
		t.Errorf("Original RDD was modified: got %v, want %v", originalResult, expectedOriginal)
	}

	// Mapped should have doubled values
	mappedResult := mapped.GetData()
	expectedMapped := []interface{}{2, 4, 6}
	if !reflect.DeepEqual(mappedResult, expectedMapped) {
		t.Errorf("Mapped RDD failed: got %v, want %v", mappedResult, expectedMapped)
	}

	// Filtered should have filtered values
	filteredResult := filtered.GetData()
	expectedFiltered := []interface{}{2, 3}
	if !reflect.DeepEqual(filteredResult, expectedFiltered) {
		t.Errorf("Filtered RDD failed: got %v, want %v", filteredResult, expectedFiltered)
	}
}

func TestRDD_EmptyChain(t *testing.T) {
	// Test RDD with no operations
	rdd := NewKeyedRDD([]interface{}{1, 2, 3}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	result := rdd.GetData()
	expected := []interface{}{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Empty chain failed: got %v, want %v", result, expected)
	}
} 