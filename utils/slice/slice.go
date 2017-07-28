package slice

func ContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, target int) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func IndexString(slice []string, target string) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == target {
			return i
		}
	}
	return -1
}

func LastIndexString(slice []string, target string) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == target {
			return i
		}
	}
	return -1
}

func RemoveString(slice, toRemove []string) (result []string) {
	if len(slice) == 0 || len(toRemove) == 0 {
		return slice
	}
	m := make(map[string]bool)
	for _, v := range toRemove {
		m[v] = true
	}
	for _, v := range slice {
		if !m[v] {
			result = append(result, v)
		}
	}
	return
}

func RemoveInt(slice, toRemove []int) (result []int) {
	if len(slice) == 0 || len(toRemove) == 0 {
		return slice
	}
	m := make(map[int]bool)
	for _, v := range toRemove {
		m[v] = true
	}
	for _, v := range slice {
		if !m[v] {
			result = append(result, v)
		}
	}
	return
}
