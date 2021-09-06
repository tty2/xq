/*
Package domain keeps all the shared structures.
*/
package domain

// Step keeps name of tag and index set by caller.
type Step struct {
	Name  string
	Index int
}

func PathsMatch(p []Step, ph []string) bool {
	if len(p) != len(ph) {
		return false
	}

	for i := range p {
		if p[i].Name != ph[i] {
			return false
		}
	}

	return true
}
