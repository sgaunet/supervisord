//nolint:revive // Package "util" is standard for utility functions in Go projects
package util

import "slices"

// InArray returns true if the elem is in the array arr.
func InArray(elem any, arr []any) bool {
	return slices.Contains(arr, elem)
}

// HasAllElements returns true if the array arr1 contains all elements of array arr2.
func HasAllElements(arr1 []any, arr2 []any) bool {
	for _, e2 := range arr2 {
		if !InArray(e2, arr1) {
			return false
		}
	}
	return true
}

// StringArrayToInterfacArray converts []string to []any.
func StringArrayToInterfacArray(arr []string) []any {
	result := make([]any, 0, len(arr))
	for _, s := range arr {
		result = append(result, s)
	}
	return result
}

// Sub returns all the element in arr1 but not in arr2.
func Sub(arr1 []string, arr2 []string) []string {
	result := make([]string, 0)
	for _, s := range arr1 {
		if !slices.Contains(arr2, s) {
			result = append(result, s)
		}
	}
	return result
}

// IsSameStringArray returns true if arr1 and arr2 has exactly same elements without order care.
func IsSameStringArray(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for _, s := range arr1 {
		if !slices.Contains(arr2, s) {
			return false
		}
	}
	return true
}
