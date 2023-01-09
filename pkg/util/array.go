package util

//判断元素target是否在sources中存在，存在返回true, 反之返回false
func Contains[T comparable](sources []T, target T) bool {
	contains := false
	for _, s := range sources {
		if target == s {
			contains = true
			break
		}
	}
	return contains
}

//判断数组元素targets是否都在sources中存在，全部存在返回true, 反之返回false
func ContainsArray[T comparable](sources []T, targets ...T) bool {
	for _, target := range targets {
		if !Contains(sources, target) {
			return false
		}
	}
	return true
}

//取两个切片数组的交集
func Intersection[T comparable](a1, a2 []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a1 {
		hash[v] = struct{}{}
	}

	for _, v := range a2 {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}

//取两个切片数组集合的并集
func Union[T comparable](a1, a2 []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a1 {
		hash[v] = struct{}{}
		set = append(set, v)
	}

	for _, v := range a2 {
		if _, ok := hash[v]; !ok {
			hash[v] = struct{}{}
			set = append(set, v)
		}
	}

	return set
}

//取两个切片数组的diff
//返回a1无法在a2中找到的元素数组。
func Difference[T comparable](a1, a2 []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a2 {
		hash[v] = struct{}{}
	}

	for _, v := range a1 {
		if _, ok := hash[v]; !ok {
			set = append(set, v)
		}
	}
	return set
}
