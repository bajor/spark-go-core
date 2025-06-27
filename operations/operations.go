package operations

// Map applies a function to each element in a slice
func Map(data []interface{}, f func(interface{}) (interface{}, error)) ([]interface{}, error) {
	result := make([]interface{}, len(data))
	for i, item := range data {
		transformed, err := f(item)
		if err != nil {
			return nil, err
		}
		result[i] = transformed
	}
	return result, nil
}

// Filter keeps only elements that match the predicate
func Filter(data []interface{}, f func(interface{}) bool) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range data {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce applies a function to combine all elements into a single result
func Reduce(data []interface{}, f func([]interface{}) ([]interface{}, error)) ([]interface{}, error) {
	return f(data)
}

// ReduceByKey groups elements by key and applies a function to each group
func ReduceByKey(data []interface{}, keyFunc func(interface{}) (interface{}, error), reduceFunc func([]interface{}) ([]interface{}, error)) ([]interface{}, error) {
	// Group data by key
	groups := make(map[interface{}][]interface{})
	for _, item := range data {
		key, err := keyFunc(item)
		if err != nil {
			return nil, err
		}
		groups[key] = append(groups[key], item)
	}

	// Apply reduce function to each group
	result := make([]interface{}, 0)
	for _, group := range groups {
		reduced, err := reduceFunc(group)
		if err != nil {
			return nil, err
		}
		result = append(result, reduced...)
	}

	return result, nil
} 