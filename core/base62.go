package core

import (
	"fmt"
)

const (
	lowercaseStart = 0
	uppercaseStart = 26
	numericStart   = 52
)

var set = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Adapted from https://dev.to/joshduffney/what-is-base62-conversion-13o0.
func base62Encode(n uint64) []byte {
	if n == 0 {
		return []byte{set[0]}
	}
	var base62 []byte
	for n > 0 {
		remainder := n % 62
		base62 = append(base62, set[remainder])
		n /= 62
	}
	return base62
}

func base62Decode(str []byte) uint64 {
	var n uint64
	var (
		num int
	)
	for i, b := range str {
		if b >= 'a' {
			num = lowercaseStart + int(b-'a')
		} else if b >= 'A' {
			num = uppercaseStart + int(b-'A')
		} else {
			num = numericStart + int(b-'0')
		}
		v := uint64(num * pow(62, i))
		fmt.Println(v)
		n += v
	}
	return n
}

func pow(b, e int) int {
	if e == 0 {
		return 1
	}
	if e == 1 {
		return b
	}
	n := b
	for i := int(2); i <= e; i++ {
		n *= b
	}
	return n
}
