package main

var headSlice = []byte("HEAD /")
var headSliceLength = len(headSlice)

func main() {
	runner := NewRunner("tcp", "127.0.0.1:80", "big.txt")
	runner.Start()
}
