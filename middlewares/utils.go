package middlewares

func FindStringInSlice(arr []string, element string) bool {
	for _, value := range arr {
		if value == element {
			return true
		}
	}
	return false
}
