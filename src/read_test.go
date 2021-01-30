package gogasm

import (
	"bufio"
	"os"
	"testing"
	"io/ioutil"
	"strings"
)

var file string = "../server/big.txt"

func BenchmarkStreamRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open(file)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		out := make([]string, 0)
		for scanner.Scan() {
			out = append(out, scanner.Text())
		}

		file.Close()
	}
}

func BenchmarkStringsRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data, _ := ioutil.ReadFile(file)
		lines := strings.Split(string(data), "\n")

		out := make([]string, 0)

		for i := range lines {
			out = append(out, lines[i])
		}
	}
}
