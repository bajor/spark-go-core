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

func TestRDD_LazyEvaluationWithSideEffects(t *testing.T) {
	// Test lazy evaluation by using side effects to track execution
	executionCount := 0
	
	rdd := NewKeyedRDD([]interface{}{1, 2, 3}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	// Add operations that increment a counter
	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		executionCount++
		return i.(int) * 2, nil
	})

	rdd = rdd.Filter(func(i interface{}) bool {
		executionCount++
		return i.(int) > 2
	})

	// Before GetData(), no operations should have been executed
	if executionCount != 0 {
		t.Errorf("Operations were executed before GetData(): got %d executions, want 0", executionCount)
	}

	// Now call GetData() - operations should be executed
	_ = rdd.GetData()

	// Should have executed 6 operations (3 for map + 3 for filter)
	if executionCount != 6 {
		t.Errorf("Wrong number of operations executed: got %d, want 6", executionCount)
	}
}

func TestRDD_MultipleGetDataCalls(t *testing.T) {
	// Test that multiple GetData() calls re-execute operations
	executionCount := 0
	
	rdd := NewKeyedRDD([]interface{}{1, 2, 3}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		executionCount++
		return i.(int) * 2, nil
	})

	// First call
	result1 := rdd.GetData()
	expected := []interface{}{2, 4, 6}
	if !reflect.DeepEqual(result1, expected) {
		t.Errorf("First GetData() failed: got %v, want %v", result1, expected)
	}
	firstExecutionCount := executionCount

	// Second call - should re-execute operations
	result2 := rdd.GetData()
	if !reflect.DeepEqual(result2, expected) {
		t.Errorf("Second GetData() failed: got %v, want %v", result2, expected)
	}
	
	// Should have executed operations twice
	if executionCount != firstExecutionCount*2 {
		t.Errorf("Operations not re-executed: got %d total executions, want %d", executionCount, firstExecutionCount*2)
	}
}

func TestRDD_LazyEvaluationWithComplexChain(t *testing.T) {
	// Test lazy evaluation with a complex chain of operations
	executionCount := 0
	
	rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	// Build a complex chain
	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		executionCount++
		return i.(int) * 2, nil
	})

	rdd = rdd.Filter(func(i interface{}) bool {
		executionCount++
		return i.(int) > 4
	})

	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		executionCount++
		return i.(int) + 1, nil
	})

	rdd = rdd.Filter(func(i interface{}) bool {
		executionCount++
		return i.(int) > 7
	})

	// Before evaluation
	if executionCount != 0 {
		t.Errorf("Operations executed before GetData(): got %d, want 0", executionCount)
	}

	// Evaluate
	result := rdd.GetData()
	expected := []interface{}{9, 11, 13} // 9, 11, and 13 are > 7 after all operations
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Complex lazy evaluation failed: got %v, want %v", result, expected)
	}

	// Should have executed operations for all elements that passed through the chain
	// 6 elements * 2 operations + 4 elements * 2 operations = 20 total executions
	if executionCount != 20 {
		t.Errorf("Wrong execution count: got %d, want 20", executionCount)
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

func TestRDD_IntegrationWithOperations(t *testing.T) {
	// Integration test: verify RDD correctly uses operations package
	rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	// Test that RDD operations produce same results as direct operations
	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	}).Filter(func(i interface{}) bool {
		return i.(int) > 4
	})

	result := rdd.GetData()
	expected := []interface{}{6, 8, 10, 12}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RDD integration failed: got %v, want %v", result, expected)
	}
} 