package internal

import "math/rand"

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	urlPrefix    = "http://example.com/"
	urlLength    = 10
	PostEndpoint = "http://localhost:1323/shorten?url="
)

func GenerateRandomUrl() string {
	randomString := make([]byte, urlLength)
	for i := range randomString {
		randomString[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return urlPrefix + string(randomString)
}
