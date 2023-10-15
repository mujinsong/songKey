package myUtils

func IsEmpty(str string) bool {
	if &str == nil {
		return true
	}
	if len(str) == 0 {
		return true
	}
	if str == "" {
		return true
	}
	return false
}
