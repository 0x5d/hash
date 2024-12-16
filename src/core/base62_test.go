package core

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func TestBase62EncodeDecode(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Run(fmt.Sprintf("Encode/Decode test %d", i), func(t *testing.T) {
			n := rand.Uint64()
			encoded := base62Encode(n)
			decoded, err := base62Decode(encoded)
			if err != nil {
				t.Fatalf("unexpected err: %s", err)
			}
			if decoded != n {
				t.Fatalf("Expected %q to decode into %d, but got %d", encoded, n, decoded)
			}
		})
	}
}
