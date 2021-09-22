package autocoins

import "sort"

func ContainsString(a []string, x string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func ContainsStringSorted(a []string, x string) bool {
	n := sort.SearchStrings(a, x)
	return n < len(a) && a[n] == x
}
