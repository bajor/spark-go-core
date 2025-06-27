The README will be structured as follows: I will outline all the features that have been implemented, along with code examples to illustrate their functionality.

## Lazy evaluation
The implementation includes chaining of map, filter, and reduce operations. There is simple execution optimization:

- First it filters out the data.
- Then mapping operations are applied.
- Finally, the data is reduced.

### Filter Behavior
The filter operation works as follows:
- If an element matches the filter condition, it is left as is
- If an element doesn't match the filter condition, it is converted in-place to `nil`
- All other operations (map, reduce) skip `nil` elements
- In the final evaluation, all `nil` elements are removed before returning the result

This approach allows for deferred filtering where non-matching elements are marked as `nil` and skipped by subsequent operations, with final cleanup happening during evaluation.

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

The RDD (Resilient Distributed Dataset) implementation provides a high-level API for data processing with lazy evaluation.

**Code example:**
```go
package main

import (
	"fmt"

	RDD "github.com/bajor/spark-go-core/resiliant_distributed_dataset"
)

func main() {
	// Create a new RDD with data [1, 2, 3, 4, 5, 6]
	rdd := RDD.NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	// Map: multiply each element by 2
	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	})

	// Filter: keep only elements > 4
	rdd = rdd.Filter(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	})

	// Get the result: [10, 12] (only 10 and 12 are > 4)
	fmt.Println("After filter:", rdd.GetData())

	// Reduce: count elements
	rdd = rdd.ReduceByKey(func(a []interface{}) ([]interface{}, error) {
		return []interface{}{len(a)}, nil
	})

	fmt.Println("Count:", rdd.GetData())
}
```

TODO
#### first iteration
Any transformation on an RDD creates a new RDD

Data Partitioning: Decide on a strategy for data partitioning across nodes. This could be based on data size or a key in the data itself (key-based partitioning).

#### next iteraions 
Distributed Computation: Use goroutines for concurrent execution of tasks. Channels can be used for communication between these tasks, especially for shuffling data during transformations like reduceByKey.

Communication: Implement a networking layer to allow nodes to communicate. Go's standard net package can be used for TCP/UDP communication.
Serialization: Data sent over the network needs to be serialized. Go supports several serialization formats (JSON, Protobuf, etc.) that you can leverage for efficient data transfer.