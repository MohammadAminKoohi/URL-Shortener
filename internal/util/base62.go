package util

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func ToBase62(url string) string {
	urlBytes := []byte(url)

	var num uint64
	for _, b := range urlBytes {
		num = num*256 + uint64(b)
	}

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
