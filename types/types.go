package types

// KeyedRDD represents a Resilient Distributed Dataset with key-based operations
type KeyedRDD struct {
	Data  []interface{}
	Chain *OperationChain
	Key   func(i interface{}) (interface{}, error)
}

// OperationChain holds a sequence of operations to be executed lazily
type OperationChain struct {
	Operations []Operation
}

// Operation interface defines the contract for all RDD operations
type Operation interface {
	Execute(data []interface{}) ([]interface{}, error)
} 