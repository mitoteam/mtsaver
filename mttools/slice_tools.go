package mttools

// Returns number of `value` values found in `values_list`
func CountValues(value interface{}, values_list ...interface{}) (count int) {
	count = 0

	for _, element := range values_list {
		if element == value {
			count++
		}
	}
	return count
}
