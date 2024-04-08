The README will be structured as follows: I will outline all the features that have been implemented, along with code examples to illustrate their functionality.

## Lazy evaluation
The implementation includes chaining of map, filter, and reduce operations. There is simple execution optimization:

- First it filters out the data.
- Then mapping operations are applied.
- Finally, the data is reduced.

**Code example:**
```go
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
```

## RDDs

TODO
