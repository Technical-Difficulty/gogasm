package gogasm

import (
	"testing"
)

// https://stackoverflow.com/questions/41460750/how-to-convert-utf8-string-to-byte

var result interface{}

func BenchmarkStringToBytesConversion(b *testing.B) {
	s := "HEAD / HTTP/1.0\r\n\r\n"
	var bytes []byte
	b.ResetTimer()
	for n := 0; n <= b.N; n++ {
		bytes = []byte(s)
	}
	b.StopTimer()

	result = bytes
}

func BenchmarkStringToBytesByteByByte(b *testing.B) {
	s := "HEAD / HTTP/1.0\r\n\r\n"
	var bytes []byte
	b.ResetTimer()
	for n := 0; n <= b.N; n++ {
		for i := 0; i < len(s); i++ {
			bytes = append(bytes, s[i])
		}

		bytes = nil
	}
	b.StopTimer()
	result = bytes
}

// This is lightning fast if the []byte make() is called outside of the loop, it's
// also reusable so we could perhaps pre-allocate the byte slices for various
// string lengths to cut down on mem allocation.
func BenchmarkCopyStringToBytes(b *testing.B) {
	s := "HEAD / HTTP/1.0\r\n\r\n"
	bytes := make([]byte, len(s))
	b.ResetTimer()
	for n := 0; n <= b.N; n++ {
		copy(bytes[:], s)
	}
}

// Can't use fixed array for conn.Write but interesting to see the speed
func BenchmarkCopyStringToArray(b *testing.B) {
	s := "HEAD / HTTP/1.0\r\n\r\n"
	bytes := [19]byte{}
	b.ResetTimer()
	for n := 0; n <= b.N; n++ {
		copy(bytes[:], s)
	}
}
