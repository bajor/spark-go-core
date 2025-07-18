The README will be structured as follows: I will outline all the features that have been implemented, along with code examples to illustrate their functionality.

## Lazy evaluation

The implementation supports chaining of map, filter, and reduce operations. Operations are only evaluated when you call `GetData()`.

**Example:**
```go
rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) { return i, nil })

// Chain operations (no evaluation yet)
rdd = rdd.Map(func(i interface{}) (interface{}, error) { return i.(int) * 2, nil })
rdd = rdd.Filter(func(i interface{}) bool { return i.(int) > 4 })

// Evaluation happens here
result := rdd.GetData() // [6 8 10 12]
```

## RDD Chaining and Reduce

```go
rdd := NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) { return i, nil })

// Chain transformations
rdd = rdd.Map(func(i interface{}) (interface{}, error) { return i.(int) * 2, nil })
rdd = rdd.Filter(func(i interface{}) bool { return i.(int) > 4 })

// Reduce by key (count elements)
rdd = rdd.ReduceByKey(func(a []interface{}) ([]interface{}, error) { return []interface{}{len(a)}, nil })

result := rdd.GetData() // [1 1 1 1]
```

## RDD Immutability

```go
original := NewKeyedRDD([]interface{}{1, 2, 3}, func(i interface{}) (interface{}, error) { return i, nil })

mapped := original.Map(func(i interface{}) (interface{}, error) { return i.(int) * 2, nil })
filtered := original.Filter(func(i interface{}) bool { return i.(int) > 1 })

original.GetData() // [1 2 3]
mapped.GetData()   // [2 4 6]
filtered.GetData() // [2 3]
```

## TODO

### Simple Distributed POC Implementation

#### Phase 1: Basic Node Emulation
- [ ] **Simple Driver Node**
  - TCP server listening for worker connections
  - Basic RDD operation distribution
  - Simple worker registration (in-memory list)
- [ ] **Simple Worker Nodes**
  - TCP client connecting to driver
  - Basic task execution (map/filter operations)
  - Return results to driver
- [ ] **Basic Communication**
  - Simple JSON messages over TCP
  - Basic request/response pattern
  - No complex error handling

#### Phase 2: RDD Partitioning
- [ ] **Simple Data Partitioning**
  - Split RDD data across available workers
  - Round-robin or hash-based distribution
  - Basic partition metadata tracking
- [ ] **Distributed RDD Operations**
  - Send map/filter operations to workers
  - Collect results from all workers
  - Simple result aggregation

#### Phase 3: Proof of Concept
- [ ] **Basic Task Distribution**
  - Driver splits RDD into partitions
  - Send partitions + operations to workers
  - Collect and combine results
- [ ] **Simple Execution**
  - One operation at a time (no complex chaining)
  - Basic result collection
  - Simple error handling (panic on failure)

#### Phase 4: Demo Setup
- [ ] **Manual Node Startup**
  - Start driver with configurable port
  - Start multiple worker processes manually
  - Simple configuration via command line flags
- [ ] **Basic Demo**
  - Simple map/filter operations across nodes
  - Show data partitioning and result collection
  - Basic logging to prove distribution works
