/*
Package slice keeps several helpers for slice structure.
*/
package slice

// ContainsString if slice `s` contains value `v`.
func ContainsString(s []string, v string) bool {
	for i := range s {
		if s[i] == v {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicates from slice.
func RemoveDuplicates(s []string) []string {
	res := make([]string, 0, len(s))

	for i := range s {
		if ContainsString(res, s[i]) {
			continue
		}
		res = append(res, s[i])
	}

	return res
}
