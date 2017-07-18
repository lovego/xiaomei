package union

func UnionAndSort(a, b interface{}, unionKey string, sortKey string) [][2]interface{} {
	data := Union(a, b, unionKey)
	UnionSort(data, sortKey)
	return data
}
