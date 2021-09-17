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
