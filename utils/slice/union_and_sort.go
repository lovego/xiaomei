package slice

func UnionAndSort(a, b interface{}, union_key string, sort_key string) [][2]interface{} {
	data := Union(a, b, union_key)
	UnionSort(data, sort_key)
	return data
}
