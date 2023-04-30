package functions

import "fmt"

func Equal(b [10]byte, s string) bool {
	if len(b) < len(s) {
		fmt.Println("la cadena es mas larga que el array de bytes")
		return false
	}
	for i, x := range s {
		if byte(x) != b[i] {
			return false
		}
	}
	return true
}

func TrimArray(arr []byte) []byte {
	var result []byte
	for _, v := range arr {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}
