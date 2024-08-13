package util

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func ToBase62(url string) string {
	// Step 1: Convert the URL string to a byte slice
	urlBytes := []byte(url)

	// Step 2: Iterate over the byte slice to create a uint64 number
	var num uint64
	for _, b := range urlBytes {
		num = num*256 + uint64(b)
	}

	// Step 3: Convert the number to a Base62 string
	if num == 0 {
		return string(base62Chars[0])
	}
	var result []byte
	for num > 0 {
		remainder := num % 62
		result = append([]byte{base62Chars[remainder]}, result...)
		num = num / 62
	}
	return string(result)
}
