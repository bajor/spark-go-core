package lazy

import "testing"

func TestLazyChainCollectReset(t *testing.T) {
	lc := NewLazyChain([]interface{}{1, 2, 3})
	lc.AddMap(func(i interface{}) (interface{}, error) { return i.(int) * 2, nil })

	first, err := lc.Collect()
	if err != nil {
		t.Fatalf("first collect error: %v", err)
	}
	second, err := lc.Collect()
	if err != nil {
		t.Fatalf("second collect error: %v", err)
	}

	if len(first) != len(second) {
		t.Fatalf("expected same length, got %d and %d", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("results differ: %v vs %v", first, second)
		}
	}
}

func TestLazyChainForEachReset(t *testing.T) {
	lc := NewLazyChain([]interface{}{1, 2, 3})
	sum1 := 0
	if err := lc.ForEach(func(i interface{}) error {
		sum1 += i.(int)
		return nil
	}); err != nil {
		t.Fatalf("first ForEach error: %v", err)
	}
	if sum1 != 6 {
		t.Fatalf("expected sum 6, got %d", sum1)
	}

	sum2 := 0
	if err := lc.ForEach(func(i interface{}) error {
		sum2 += i.(int)
		return nil
	}); err != nil {
		t.Fatalf("second ForEach error: %v", err)
	}
	if sum2 != 6 {
		t.Fatalf("expected sum 6 again, got %d", sum2)
	}
}
