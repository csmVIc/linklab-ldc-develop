package tool

// 判断Map A中是否包含Map B
func MapAIncludeMapB(a map[string]string, b map[string]string) bool {
	for bKey, bValue := range b {
		if aValue, isOK := a[bKey]; isOK == false {
			return false
		} else if aValue != bValue {
			return false
		}
	}
	return true
}
