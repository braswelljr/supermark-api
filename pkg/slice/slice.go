package slice

// Contains - return true/false if an element is in a slice or not
//
//	@param slice - slice to check
//	@param value - value to check for
//	@return bool - true/false
func Contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// Remove element by value from slice
//
//	@param l - slice to remove from
//	@param item - item to remove
//	@return []T - slice without the item
func Remove[T comparable](l []T, item T) []T {
	for i, ele := range l {
		if ele == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// Filter - filter slice by conditions
//
//	@param slice - slice to filter
//	@param conditions - conditions to filter by (array of functions)
//	@return []T - filtered slice
func Filter[T comparable](slice []T, conditions ...func(T) bool) []T {
	// results for the condition
	var result []T

	// if the slice or the conditions are empty
	if len(slice) == 0 || len(conditions) == 0 {
		// return the original slice
		return slice
	}

	// loop through the slice
	for _, item := range slice {
		// match flag
		match := true
		// loop through the conditions
		for _, condition := range conditions {
			// if the condition is not met
			if !condition(item) {
				// set the match flag to false and
				match = false

				// break the loop
				break
			}
		}

		// if the match flag is true
		if match {
			// append the item to the results
			result = append(result, item)
		}
	}
	// if the result is empty
	if len(result) == 0 {
		// return the original slice
		result = slice
	}

	// return the result
	return result
}
