package rdd

// KeyedRDD represents a Resilient Distributed Dataset with key-based operations
type KeyedRDD struct {
	data  []interface{}
	chain *OperationChain
	key   func(i interface{}) (interface{}, error)
}

// OperationChain holds a sequence of operations to be executed lazily
type OperationChain struct {
	operations []Operation
}

// Operation interface defines the contract for all RDD operations
type Operation interface {
	Execute(data []interface{}) ([]interface{}, error)
} 