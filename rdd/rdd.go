package rdd

type rddOpType int

const (
	rddMapOp rddOpType = iota
	rddFilterOp
	rddReduceByKeyOp
	rddReduceOp
)

type rddOp struct {
	type_ rddOpType
	mapFn func(interface{}) (interface{}, error)
	filterFn func(interface{}) bool
	reduceFn func([]interface{}) ([]interface{}, error)
}

type KeyedRDD struct {
	data  []interface{}
	ops   []rddOp
	key   func(i interface{}) (interface{}, error)
}

func NewKeyedRDD(data []interface{}, key func(i interface{}) (interface{}, error)) *KeyedRDD {
	return &KeyedRDD{
		data: data,
		ops:  nil,
		key:  key,
	}
}

func (r *KeyedRDD) Map(f func(i interface{}) (interface{}, error)) *KeyedRDD {
	newOps := append([]rddOp{}, r.ops...)
	newOps = append(newOps, rddOp{type_: rddMapOp, mapFn: f})
	return &KeyedRDD{data: r.data, ops: newOps, key: r.key}
}

func (r *KeyedRDD) Filter(f func(i interface{}) bool) *KeyedRDD {
	newOps := append([]rddOp{}, r.ops...)
	newOps = append(newOps, rddOp{type_: rddFilterOp, filterFn: f})
	return &KeyedRDD{data: r.data, ops: newOps, key: r.key}
}

func (r *KeyedRDD) ReduceByKey(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newOps := append([]rddOp{}, r.ops...)
	newOps = append(newOps, rddOp{type_: rddReduceByKeyOp, reduceFn: f})
	return &KeyedRDD{data: r.data, ops: newOps, key: r.key}
}

func (r *KeyedRDD) Reduce(f func(a []interface{}) ([]interface{}, error)) *KeyedRDD {
	newOps := append([]rddOp{}, r.ops...)
	newOps = append(newOps, rddOp{type_: rddReduceOp, reduceFn: f})
	return &KeyedRDD{data: r.data, ops: newOps, key: r.key}
}

func (r *KeyedRDD) GetData() []interface{} {
	data := r.data
	for _, op := range r.ops {
		switch op.type_ {
		case rddMapOp:
			var out []interface{}
			for _, v := range data {
				res, err := op.mapFn(v)
				if err != nil {
					panic(err)
				}
				out = append(out, res)
			}
			data = out
		case rddFilterOp:
			var out []interface{}
			for _, v := range data {
				if op.filterFn(v) {
					out = append(out, v)
				}
			}
			data = out
		case rddReduceByKeyOp:
			partitioned := make(map[interface{}][]interface{})
			for _, item := range data {
				key, err := r.key(item)
				if err != nil {
					panic(err)
				}
				partitioned[key] = append(partitioned[key], item)
			}
			var out []interface{}
			for _, part := range partitioned {
				res, err := op.reduceFn(part)
				if err != nil {
					panic(err)
				}
				out = append(out, res...)
			}
			data = out
		case rddReduceOp:
			res, err := op.reduceFn(data)
			if err != nil {
				panic(err)
			}
			data = res
		}
	}
	return data
} 