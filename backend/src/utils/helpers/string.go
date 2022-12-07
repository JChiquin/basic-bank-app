package helpers

/*
StringInSlice receives a string and a slice
then tries to find the string in the list
*/
func StringInSlice(toFind string, list []string) bool {
	for _, validItems := range list {
		if toFind == validItems {
			return true
		}
	}
	return false
}

/*
Custom type so Request.Context() doesn't complain
*/
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}
