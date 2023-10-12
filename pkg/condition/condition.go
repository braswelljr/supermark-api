package condition

// Ternary - is a function that returns a value based on a condition.
//
//	@param condition - bool
//	@param a - interface{}
//	@param b - interface{}
//	@return - interface{}
func Ternary[T any](condition bool, a T, b T) T {
	// check if the condition is true
	if condition {
		return a
	}

	// return b
	return b
}
