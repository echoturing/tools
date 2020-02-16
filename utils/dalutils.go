package utils

func WrapWith(s []string, wrapper string) []string {
	result := make([]string, 0, len(s))
	for i := range s {
		result = append(result, wrapper+s[i]+wrapper)
	}
	return result
}

// StringJoin
// TODO(xiangxu) improve performance
func StringJoin(s []string, joiner string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[1]
	}
	result := ""
	for i := range s {
		if i == 0 {
			result += s[i]
		} else {
			result += joiner + s[i]
		}
	}
	return result
}
