package strings

import "strings"

//SnackToCamel snack string to camel
func SnackToCamel(str string) string {
	strArr := strings.Split(str, "_")
	newArr := make([]string, 0, len(strArr))
	for _, v := range strArr {
		newArr = append(newArr, strings.Title(v))
	}
	return strings.Join(newArr, "")
}
