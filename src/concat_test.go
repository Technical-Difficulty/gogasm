package gogasm

import (
	"bytes"
	"strings"
	"testing"
)

var (
	headSlice = []byte{72, 69, 65, 68, 32, 47}
	testSlice = []byte{116, 101, 115, 116, 105, 110, 103, 49, 50, 51}
	httpSlice = []byte{32, 72, 84, 84, 80, 47, 49, 46, 48, 13, 10, 13, 10}
)

// https://stackoverflow.com/questions/32370615/what-is-the-fastest-way-to-concatenate-several-byte-together?noredirect=1
// https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go

func BenchmarkStringConcat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var str string
		str += "x"
	}
}

func BenchmarkStringBuffer(b *testing.B) {
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")
	}
}

//func BenchmarkByteSliceBuffer(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		var buffer bytes.Buffer
//		buffer.wr
//	}
//}

func BenchmarkStringCopy(b *testing.B) {
	bs := make([]byte, b.N)
	bl := 0

	for n := 0; n < b.N; n++ {
		bl += copy(bs[bl:], "x")
	}
}

// Go 1.10
func BenchmarkStringBuilder(b *testing.B) {
	var strBuilder strings.Builder

	for n := 0; n < b.N; n++ {
		strBuilder.WriteString("x")
	}
}

func BenchmarkByteSliceCopyPreAllocate(b *testing.B) {
	bytes := make([]byte, len(headSlice)+len(testSlice)+len(httpSlice))
	copy(bytes[0:len(headSlice)], headSlice)
	copy(bytes[len(headSlice)+len(testSlice):], httpSlice)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		copy(bytes[len(headSlice):], testSlice)
	}
}

func BenchmarkStringBuilderByteConvert(b *testing.B) {
	var strBuilder strings.Builder

	for n := 0; n < b.N; n++ {
		strBuilder.WriteString("HEAD /")
		strBuilder.WriteString("testing123")
		strBuilder.WriteString(" HTTP/1.0\r\n\r\n")

		result = []byte(strBuilder.String())
		strBuilder.Reset()
	}
}
