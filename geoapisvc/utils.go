package main

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func validStrId(id string) bool {
	return !strings.ContainsAny(id, "`~!@#$%^&*()+={[}]|\\:;\"'<,>.?/")
}

func inIntSlice(item int, slide []int) bool {
	for _, i := range slide {
		if item == i {
			return true
		}
	}
	return false
}
