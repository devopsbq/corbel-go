package silkroad

// stringInSlice looks if a string is in a string array
func stringInSlice(array []string, item string) bool {
	for _, i := range array {
		if i == item {
			return true
		}
	}
	return false
}
