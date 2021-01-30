package main

import (
	"fmt"
	"flag"
	gg "gogasm/src"
)

type config struct {
	wordlist string
	addr string
	port string
}

func parseFlags() *config {
	c := &config{}
	flag.StringVar(&c.wordlist, "w", "", "Wordlist file path")
	flag.StringVar(&c.addr, "a", "", "Address eg. localhost")
	flag.StringVar(&c.port, "p", "80", "Port eg. 80")
	flag.Parse()
	
	return c
}

func main() {
	config := parseFlags()

	if config.wordlist == "" {
		fmt.Println("Please provide a wordlist :: $ gogasm -w /path/to/wordlist")
		return
	}

	if config.addr == "" {
		fmt.Println("Please provide an address :: $ gogasm -a localhost")
		return 
	}

	addr := fmt.Sprintf("%s:%s", config.addr, config.port)
	runner := gg.NewRunner("tcp", addr, config.wordlist)
	runner.Start()
}
