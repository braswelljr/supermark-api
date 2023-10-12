package slice

import (
	"testing"
)

// TestContains - test the Contains function
//
//	@param t - testing.T
func TestContains(t *testing.T) {
	// create a slice
	slice := []struct {
		slice      []string
		contains   string
		itContains bool
	}{
		{
			slice:      []string{"a", "b", "c"},
			contains:   "a",
			itContains: true,
		},
		{
			slice:      []string{"a", "b", "c"},
			contains:   "d",
			itContains: false,
		},
	}

	// check if the slice contains the value
	for _, item := range slice {
		// check if the slice contains the value
		if Contains(item.slice, item.contains) != item.itContains {
			// log the error
			t.Errorf("slice %v should contain %v", item.slice, item.contains)
		}
	}
}

// TestRemove - test the Remove function
//
//	@param t - testing.T
func TestRemove(t *testing.T) {
	// create a slice
	slice := []struct {
		slice    []string
		contains string
	}{
		{
			slice:    []string{"a", "b", "c"},
			contains: "a",
		},
		{
			slice:    []string{"a", "b", "c"},
			contains: "d",
		},
	}

	// check if the slice contains the value
	for _, item := range slice {
		// remove the item from the slice
		slice := Remove(item.slice, item.contains)

		// check if the slice contains the value
		if Contains(slice, item.contains) {
			// log the error
			t.Errorf("slice %v should not contain %v", slice, item.contains)
		}
	}
}

// TestFilter - test the Filter function
//
//	@param t - testing.T
func TestFilter(t *testing.T) {
	// create a slice
	slice := []struct {
		slice      []string
		contains   string
		itContains bool
	}{
		{
			slice:      []string{"a", "b", "c"},
			contains:   "a",
			itContains: true,
		},
		{
			slice:      []string{"a", "b", "c"},
			contains:   "d",
			itContains: false,
		},
	}

	// check if the slice contains the value
	for _, item := range slice {
		// filter the slice
		slice := Filter(item.slice, func(s string) bool {
			return s == item.contains
		})

		// check if the slice contains the value
		if Contains(slice, item.contains) != item.itContains {
			// log the error
			t.Errorf("slice %v should contain %v", slice, item.contains)
		}
	}
}
