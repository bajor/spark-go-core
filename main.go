package main

import (
	"fmt"

	RDD "github.com/bajor/spark-go-core/resiliant_distributed_dataset"
)

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

	fmt.Println("After filter:", rdd.GetData())

	// Create a new RDD for the reduce operations
	rdd2 := RDD.NewKeyedRDD(rdd.GetData(), func(i interface{}) (interface{}, error) {
		return i, nil
	})

	rdd2 = rdd2.ReduceByKey(func(a []interface{}) ([]interface{}, error) {
		fmt.Println("ReduceByKey input:", a)
		return []interface{}{len(a)}, nil
	})

	fmt.Println("After ReduceByKey:", rdd2.GetData())

	rdd3 := RDD.NewKeyedRDD(rdd2.GetData(), func(i interface{}) (interface{}, error) {
		return i, nil
	})

	rdd3 = rdd3.Reduce(func(a []interface{}) ([]interface{}, error) {
		fmt.Println("Reduce input:", a)
		sum := 0
		for _, val := range a {
			sum += val.(int)
		}
		return []interface{}{sum}, nil
	})

	result := rdd3.GetData()
	fmt.Println("Final result:", result)
}
