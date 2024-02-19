package lazy

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
func (lc *LazyChain) Evaluate(inputs []interface{}) (interface{}, error) {
	var err error

	for i := 0; i < len(inputs); i++ {
		for _, op := range lc.operations {
			// for _, current := range inputs {
			inputs[i], err = op(inputs[i])
			// current, err = op(current)
			if err != nil {
				return nil, err // Early return on first error
			}
		}
	}
	return inputs, nil
}
