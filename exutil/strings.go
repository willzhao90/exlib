package exutil

func StringSliceContains(a []interface{}, x string) bool {
	for _, n := range a {
		if s, ok := n.(string); ok && x == s {
			return true
		}
	}
	return false
}
