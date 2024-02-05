package main

import "fmt"

// Operation defines a function that takes an input and returns an output with an error
type Operation func(interface{}) (interface{}, error)

// LazyChain holds a slice of operations to be executed lazily
type LazyChain struct {
    operations []Operation
}

// Add adds a new operation to the chain
func (lc *LazyChain) Add(op Operation) {
    lc.operations = append(lc.operations, op)
}

// Evaluate executes the chain of operations lazily
func (lc *LazyChain) Evaluate(input interface{}) (interface{}, error) {
    var err error
    current := input
    for _, op := range lc.operations {
        current, err = op(current)
        if err != nil {
            return nil, err // Early return on first error
        }
    }
    return current, nil
}

func main() {
    // Example usage
    chain := &LazyChain{}

    // Add operations to the chain
    chain.Add(func(i interface{}) (interface{}, error) {
        // Example operation: add 5
        return i.(int) + 5, nil
    })
    chain.Add(func(i interface{}) (interface{}, error) {
        // Another operation: multiply by 2
        return i.(int) * 2, nil
    })

    // Evaluate the chain
    result, err := chain.Evaluate(10) // Starting input is 10
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
