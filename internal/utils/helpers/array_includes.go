package helpers

func IncludesString(item string, array []string) bool {
	for _, k := range array {
		if k == item {
			return true
		}
	}
	return false
}
