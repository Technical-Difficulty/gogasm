package main

import gg "gogasm/src"

func main() {
	runner := gg.NewRunner("tcp", "127.0.0.1:80", "server/big.txt")
	runner.Start()
}
