package main

import (
	"fmt"

	RDD "github.com/bajor/spark-go-core/resiliant_distributed_dataset"
)

/*

TODO FILTER, ITS BROKE

it should work as follows:
	- if you match condition - leave it as is 
	- if you don't - convert it inplace to nil
	- all other methods should skip element if it's nil
	- in evaluate, reduce - nil should be removed

add this to readme 
*/


func main() {
	rdd := RDD.NewKeyedRDD([]interface{}{1, 2, 3, 4, 5, 6}, func(i interface{}) (interface{}, error) {
		return i, nil
	})

	rdd = rdd.Map(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	})

	rdd = rdd.Filter(func(i interface{}) bool {
		if val, ok := i.(int); ok {
			return val > 4
		}
		return false
	})

	fmt.Println("Result:", rdd.GetData())

	rdd = rdd.ReduceByKey(func(a []interface{}) ([]interface{}, error) {
		return []interface{}{len(a)}, nil
	})

	fmt.Println("Result:", rdd.GetData())

	rdd = rdd.Reduce(func(a []interface{}) ([]interface{}, error) {
		sum := 0
		for _, val := range a {
			sum += val.(int)
		}
		return []interface{}{sum}, nil
	})

	result := rdd.GetData()
	fmt.Println("Result:", result)
}
